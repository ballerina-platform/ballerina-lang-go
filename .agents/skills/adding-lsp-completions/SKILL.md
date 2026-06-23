---
name: adding-lsp-completion
description: Adding a completion for lsp given some pattern
---

First check if we already have a completion context that matches the pattern or one that gives the same completions. If there are such candidate context ask user if you should reuse it. If not you need to first setup a completion context.
1. First for given pattern create some .bal sources and see what is are the st nodes you are getting.
2. Then check what are the AST nodes you get in recovering mode.
  - You should get enough information to properly represent the pattern using AST nodes when going from ST to AST with recovering nodes. If not try to figure out a way to fix this an prompt user with that fix. State the ST nodes you are getting and AST nodes after applying the fix.
  -  The validate your fix to node builder is working properly
3. Then define a completion context and add a way to match that using AST nodes you are getting
