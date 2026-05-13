// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"ballerina-lang-go/pal"
	"ballerina-lang-go/projects"
	"ballerina-lang-go/tools/diagnostics"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type nativeHTTPClient struct {
	client *http.Client
}

func (c *nativeHTTPClient) Execute(method, url string, body []byte, contentType string, reqHeaders map[string][]string) (int, map[string][]string, []byte, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return 0, nil, nil, err
	}
	req.Header.Set("User-Agent", "ballerina")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for k, vals := range reqHeaders {
		if len(vals) == 0 {
			continue
		}
		req.Header.Set(k, vals[0])
		for _, v := range vals[1:] {
			req.Header.Add(k, v)
		}
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, nil, nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	respBody, err := io.ReadAll(resp.Body)
	return resp.StatusCode, map[string][]string(resp.Header), respBody, err
}

var nativePal = pal.Platform{
	IO: pal.IO{
		Stdout: func(p []byte) (n int, err error) {
			return os.Stdout.Write(p)
		},
		Stderr: func(p []byte) (n int, err error) {
			return os.Stderr.Write(p)
		},
	},
	HTTP: pal.HTTP{
		NewClient: func(cfg pal.ClientConfig) pal.HTTPClient {
			tlsConfig := &tls.Config{InsecureSkipVerify: cfg.TLS.InsecureSkipVerify} //nolint:gosec
			if len(cfg.TLS.CACertPEM) > 0 {
				pool := x509.NewCertPool()
				pool.AppendCertsFromPEM(cfg.TLS.CACertPEM)
				tlsConfig.RootCAs = pool
				if !cfg.TLS.InsecureSkipVerify {
					// Go 1.15+ requires SANs for hostname verification; many self-signed and
					// Java-issued certs only set the CN field. When a custom CA is provided
					// we do our own verification so CN-only certs are accepted as a fallback.
					tlsConfig.InsecureSkipVerify = true //nolint:gosec
					tlsConfig.VerifyConnection = tlsVerifyConnectionWithCNFallback(pool)
				}
			}
			if len(cfg.TLS.ClientCertPEM) > 0 && len(cfg.TLS.ClientKeyPEM) > 0 {
				if cert, err := tls.X509KeyPair(cfg.TLS.ClientCertPEM, cfg.TLS.ClientKeyPEM); err == nil {
					tlsConfig.Certificates = []tls.Certificate{cert}
				}
			}
			transport := &http.Transport{TLSClientConfig: tlsConfig}
			protocols := new(http.Protocols)
			if cfg.HTTPVersion == "2.0" {
				protocols.SetHTTP2(true)
				protocols.SetUnencryptedHTTP2(true)
			} else {
				protocols.SetHTTP1(true)
			}
			transport.Protocols = protocols
			c := &http.Client{Timeout: cfg.Timeout, Transport: transport}
			if !cfg.FollowRedirects {
				c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				}
			}
			return &nativeHTTPClient{client: c}
		},
	},
}

// printError prints an error message in the standard Ballerina CLI format to stderr.
func printError(err error, usage string, showHelp bool) {
	printErrorTo(os.Stderr, err, usage, showHelp)
}

func printRuntimeError(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
}

// printErrorTo prints an error message in the standard Ballerina CLI format to the given writer.
func printErrorTo(w io.Writer, err error, usage string, showHelp bool) {
	_, _ = fmt.Fprintf(w, "ballerina: %s\n", err.Error())
	if usage != "" {
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintln(w, "USAGE:")
		_, _ = fmt.Fprintf(w, "    %s\n", usage)
	}
	if showHelp {
		_, _ = fmt.Fprintln(w)
		_, _ = fmt.Fprintln(w, "For more information try --help")
	}
}

// validateSourceFile validates the source file argument for the 'run' command.
// Allows zero arguments (defaults to current directory in runBallerina).
func validateSourceFile(cmd *cobra.Command, args []string) error {
	// Allow zero arguments - will default to current directory "."
	// Path validation happens in directory.Load
	return nil
}

func isTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

type outputStyle struct {
	reset, red, yellow, cyan, bold string
}

func (s outputStyle) severityColor(severity diagnostics.DiagnosticSeverity) string {
	if severity == diagnostics.Warning {
		return s.yellow
	}
	return s.red
}

func outputStyleFor(noColors bool) outputStyle {
	s := outputStyle{}
	if !noColors {
		s.reset = "\033[0m"
		s.red = "\033[31m"
		s.yellow = "\033[33m"
		s.cyan = "\033[36m"
		s.bold = "\033[1m"
	}
	return s
}

type diagnosticLocation struct {
	filePath            string
	startLine, startCol int
	endLine, endCol     int
	numWidth            int
}

func buildDiagnosticLocation(filePath string, startLine, startCol, endLine, endCol int) diagnosticLocation {
	startLineNumStr := fmt.Sprintf("%d", startLine+1)
	endLineNumStr := fmt.Sprintf("%d", endLine+1)
	numWidth := len(startLineNumStr)
	if w := len(endLineNumStr); w > numWidth {
		numWidth = w
	}
	return diagnosticLocation{
		filePath:  filePath,
		startLine: startLine,
		startCol:  startCol,
		endLine:   endLine,
		endCol:    endCol,
		numWidth:  numWidth,
	}
}

func printDiagnostics(fsys fs.FS, w io.Writer, diagResult projects.DiagnosticResult, noColors bool, de *diagnostics.DiagnosticEnv) {
	for _, d := range diagResult.Diagnostics() {
		printDiagnostic(fsys, w, d, noColors, de)
	}
}

func printDiagnostic(fsys fs.FS, w io.Writer, d diagnostics.Diagnostic, noColors bool, de *diagnostics.DiagnosticEnv) {
	s := outputStyleFor(noColors)
	printDiagnosticHeader(w, s, d)

	location := d.Location()
	if diagnostics.IsLocationEmpty(location) {
		_, _ = fmt.Fprintln(w)
		return
	}
	if !diagnostics.LocationHasSource(location) {
		_, _ = fmt.Fprintf(w, "  %s-->%s %s\n\n", s.cyan, s.reset, de.FileName(location))
		return
	}

	loc := buildDiagnosticLocation(
		de.FileName(location),
		de.StartLine(location), de.StartColumn(location),
		de.EndLine(location), de.EndColumn(location),
	)
	printDiagnosticLocation(w, s, loc)
	printSourceSnippet(w, s, loc, fsys, s.severityColor(d.DiagnosticInfo().Severity()))
	_, _ = fmt.Fprintln(w)
}

func printDiagnosticHeader(w io.Writer, s outputStyle, d diagnostics.Diagnostic) {
	info := d.DiagnosticInfo()
	codeStr := ""
	if c := info.Code(); c != "" {
		codeStr = fmt.Sprintf("[%s]", c)
	}
	_, _ = fmt.Fprintf(w, "%s%s%s%s%s: %s%s%s\n",
		s.bold, s.severityColor(info.Severity()), strings.ToLower(info.Severity().String()), codeStr, s.reset,
		s.bold, d.Message(), s.reset,
	)
}

func printDiagnosticLocation(w io.Writer, s outputStyle, loc diagnosticLocation) {
	_, _ = fmt.Fprintf(w, "%*s%s-->%s %s:%d:%d\n",
		loc.numWidth, "", s.cyan, s.reset, loc.filePath, loc.startLine+1, loc.startCol+1,
	)
	if loc.filePath != "" {
		_, _ = fmt.Fprintf(w, "%*s %s|%s\n", loc.numWidth, "", s.cyan, s.reset)
	}
}

func printSourceSnippet(w io.Writer, s outputStyle, loc diagnosticLocation, fsys fs.FS, severityColor string) {
	content, err := fs.ReadFile(fsys, loc.filePath)
	if err != nil {
		return
	}
	lines := strings.Split(string(content), "\n")
	if loc.startLine >= len(lines) {
		return
	}

	for line := loc.startLine; line <= loc.endLine && line < len(lines); line++ {
		lineContent := strings.TrimSuffix(lines[line], "\r")
		lineNumStr := fmt.Sprintf("%d", line+1)

		startCol := 0
		var endCol int

		switch {
		case loc.startLine == loc.endLine:
			startCol = loc.startCol
			endCol = loc.endCol
		case line == loc.startLine:
			startCol = loc.startCol
			endCol = len(lineContent)
		case line == loc.endLine:
			startCol = 0
			endCol = loc.endCol
		default:
			startCol = 0
			endCol = len(lineContent)
		}

		var highlightLen int
		startCol, _, highlightLen = computeTrimmedCaretSpan(lineContent, startCol, endCol)

		_, _ = fmt.Fprintf(w, "%s%*s | %s%s\n", s.cyan, loc.numWidth, lineNumStr, s.reset, lineContent)
		pointer := buildPointer(lineContent, startCol, highlightLen)
		_, _ = fmt.Fprintf(w, "%*s %s| %s%s%s\n", loc.numWidth, "", s.cyan, severityColor, pointer, s.reset)
	}
}

func computeTrimmedCaretSpan(lineContent string, startCol, endCol int) (trimStartCol, trimEndCol, highlightLen int) {
	firstNonWS := -1
	for i := 0; i < len(lineContent); i++ {
		if lineContent[i] != ' ' && lineContent[i] != '\t' {
			firstNonWS = i
			break
		}
	}
	lastNonWS := len(lineContent)
	hasNonWS := firstNonWS != -1
	if hasNonWS {
		for lastNonWS > firstNonWS && (lineContent[lastNonWS-1] == ' ' || lineContent[lastNonWS-1] == '\t') {
			lastNonWS--
		}
	}
	if !hasNonWS {
		return startCol, startCol, 0
	}
	if startCol < firstNonWS {
		startCol = firstNonWS
	}
	highlightLen = endCol - startCol
	return startCol, endCol, highlightLen
}

func buildPointer(lineContent string, startCol, highlightLen int) string {
	var b strings.Builder
	for i := 0; i < startCol && i < len(lineContent); i++ {
		if lineContent[i] == '\t' {
			b.WriteByte('\t')
		} else {
			b.WriteByte(' ')
		}
	}
	for range highlightLen {
		b.WriteByte('^')
	}
	return b.String()
}

// tlsVerifyConnectionWithCNFallback returns a VerifyConnection callback that verifies the
// server's certificate chain against rootCAs and falls back to CN-based hostname matching
// when no SANs are present. Go 1.15+ disabled CN-only hostname verification (RFC 6125 §2.3),
// but many self-signed and Java-issued certificates still rely on it.
func tlsVerifyConnectionWithCNFallback(rootCAs *x509.CertPool) func(tls.ConnectionState) error {
	return func(cs tls.ConnectionState) error {
		opts := x509.VerifyOptions{
			Roots:         rootCAs,
			Intermediates: x509.NewCertPool(),
		}
		for _, cert := range cs.PeerCertificates[1:] {
			opts.Intermediates.AddCert(cert)
		}
		if _, err := cs.PeerCertificates[0].Verify(opts); err != nil {
			return err
		}
		// cs.ServerName is the SNI hostname (no port). Try SAN-based verification first;
		// fall back to CN matching for legacy certificates.
		leaf := cs.PeerCertificates[0]
		if err := leaf.VerifyHostname(cs.ServerName); err == nil {
			return nil
		}
		return tlsMatchCN(leaf.Subject.CommonName, cs.ServerName)
	}
}

// tlsMatchCN checks whether pattern (a certificate CN) matches host.
// Supports simple wildcard patterns of the form "*.example.com".
func tlsMatchCN(pattern, host string) error {
	pattern = strings.ToLower(strings.TrimSuffix(pattern, "."))
	host = strings.ToLower(strings.TrimSuffix(host, "."))
	if pattern == host {
		return nil
	}
	if strings.HasPrefix(pattern, "*.") {
		suffix := pattern[1:] // ".example.com"
		if strings.HasSuffix(host, suffix) && strings.Count(host, ".") == strings.Count(suffix, ".") {
			return nil
		}
	}
	return fmt.Errorf("x509: certificate CN %q does not match host %q", pattern, host)
}
