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

import ballerina/io;

const int K0 = 5;
const int K1 = 22;
const int K2 = 39;
const int K3 = 56;
const int K4 = 73;
const int K5 = 90;
const int K6 = 107;
const int K7 = 124;
const int K8 = 141;
const int K9 = 158;
const int K10 = 175;
const int K11 = 192;
const int K12 = 209;
const int K13 = 226;
const int K14 = 243;
const int K15 = 260;
const int K16 = 277;
const int K17 = 294;
const int K18 = 311;
const int K19 = 328;
const int K20 = 345;
const int K21 = 362;
const int K22 = 379;
const int K23 = 396;
const int K24 = 413;
const int K25 = 430;
const int K26 = 447;
const int K27 = 464;
const int K28 = 481;
const int K29 = 498;
const int K30 = 515;
const int K31 = 532;

function bigSum(
        int p0 = K0 + 0,
        int p1 = K1 + 1,
        int p2 = K2 + 2,
        int p3 = K3 + 3,
        int p4 = K4 + 4,
        int p5 = K5 + 5,
        int p6 = K6 + 6,
        int p7 = K7 + 7,
        int p8 = K8 + 8,
        int p9 = K9 + 9,
        int p10 = K10 + 10,
        int p11 = K11 + 11,
        int p12 = K12 + 12,
        int p13 = K13 + 13,
        int p14 = K14 + 14,
        int p15 = K15 + 15,
        int p16 = K16 + 16,
        int p17 = K17 + 17,
        int p18 = K18 + 18,
        int p19 = K19 + 19,
        int p20 = K20 + 20,
        int p21 = K21 + 21,
        int p22 = K22 + 22,
        int p23 = K23 + 23,
        int p24 = K24 + 24,
        int p25 = K25 + 25,
        int p26 = K26 + 26,
        int p27 = K27 + 27,
        int p28 = K28 + 28,
        int p29 = K29 + 29,
        int p30 = K30 + 30,
        int p31 = K31 + 31,
        int p32 = K0 + 32,
        int p33 = K1 + 33,
        int p34 = K2 + 34,
        int p35 = K3 + 35,
        int p36 = K4 + 36,
        int p37 = K5 + 37,
        int p38 = K6 + 38,
        int p39 = K7 + 39,
        int p40 = K8 + 40,
        int p41 = K9 + 41,
        int p42 = K10 + 42,
        int p43 = K11 + 43,
        int p44 = K12 + 44,
        int p45 = K13 + 45,
        int p46 = K14 + 46,
        int p47 = K15 + 47,
        int p48 = K16 + 48,
        int p49 = K17 + 49,
        int p50 = K18 + 50,
        int p51 = K19 + 51,
        int p52 = K20 + 52,
        int p53 = K21 + 53,
        int p54 = K22 + 54,
        int p55 = K23 + 55,
        int p56 = K24 + 56,
        int p57 = K25 + 57,
        int p58 = K26 + 58,
        int p59 = K27 + 59,
        int p60 = K28 + 60,
        int p61 = K29 + 61,
        int p62 = K30 + 62,
        int p63 = K31 + 63,
        int p64 = K0 + 64,
        int p65 = K1 + 65,
        int p66 = K2 + 66,
        int p67 = K3 + 67,
        int p68 = K4 + 68,
        int p69 = K5 + 69,
        int p70 = K6 + 70,
        int p71 = K7 + 71,
        int p72 = K8 + 72,
        int p73 = K9 + 73,
        int p74 = K10 + 74,
        int p75 = K11 + 75,
        int p76 = K12 + 76,
        int p77 = K13 + 77,
        int p78 = K14 + 78,
        int p79 = K15 + 79,
        int p80 = K16 + 80,
        int p81 = K17 + 81,
        int p82 = K18 + 82,
        int p83 = K19 + 83,
        int p84 = K20 + 84,
        int p85 = K21 + 85,
        int p86 = K22 + 86,
        int p87 = K23 + 87,
        int p88 = K24 + 88,
        int p89 = K25 + 89,
        int p90 = K26 + 90,
        int p91 = K27 + 91,
        int p92 = K28 + 92,
        int p93 = K29 + 93,
        int p94 = K30 + 94,
        int p95 = K31 + 95,
        int p96 = K0 + 96,
        int p97 = K1 + 97,
        int p98 = K2 + 98,
        int p99 = K3 + 99) returns int {
    int total = 0;
    total = (total + p0) % 10000000;
    total = (total + p1) % 10000000;
    total = (total + p2) % 10000000;
    total = (total + p3) % 10000000;
    total = (total + p4) % 10000000;
    total = (total + p5) % 10000000;
    total = (total + p6) % 10000000;
    total = (total + p7) % 10000000;
    total = (total + p8) % 10000000;
    total = (total + p9) % 10000000;
    total = (total + p10) % 10000000;
    total = (total + p11) % 10000000;
    total = (total + p12) % 10000000;
    total = (total + p13) % 10000000;
    total = (total + p14) % 10000000;
    total = (total + p15) % 10000000;
    total = (total + p16) % 10000000;
    total = (total + p17) % 10000000;
    total = (total + p18) % 10000000;
    total = (total + p19) % 10000000;
    total = (total + p20) % 10000000;
    total = (total + p21) % 10000000;
    total = (total + p22) % 10000000;
    total = (total + p23) % 10000000;
    total = (total + p24) % 10000000;
    total = (total + p25) % 10000000;
    total = (total + p26) % 10000000;
    total = (total + p27) % 10000000;
    total = (total + p28) % 10000000;
    total = (total + p29) % 10000000;
    total = (total + p30) % 10000000;
    total = (total + p31) % 10000000;
    total = (total + p32) % 10000000;
    total = (total + p33) % 10000000;
    total = (total + p34) % 10000000;
    total = (total + p35) % 10000000;
    total = (total + p36) % 10000000;
    total = (total + p37) % 10000000;
    total = (total + p38) % 10000000;
    total = (total + p39) % 10000000;
    total = (total + p40) % 10000000;
    total = (total + p41) % 10000000;
    total = (total + p42) % 10000000;
    total = (total + p43) % 10000000;
    total = (total + p44) % 10000000;
    total = (total + p45) % 10000000;
    total = (total + p46) % 10000000;
    total = (total + p47) % 10000000;
    total = (total + p48) % 10000000;
    total = (total + p49) % 10000000;
    total = (total + p50) % 10000000;
    total = (total + p51) % 10000000;
    total = (total + p52) % 10000000;
    total = (total + p53) % 10000000;
    total = (total + p54) % 10000000;
    total = (total + p55) % 10000000;
    total = (total + p56) % 10000000;
    total = (total + p57) % 10000000;
    total = (total + p58) % 10000000;
    total = (total + p59) % 10000000;
    total = (total + p60) % 10000000;
    total = (total + p61) % 10000000;
    total = (total + p62) % 10000000;
    total = (total + p63) % 10000000;
    total = (total + p64) % 10000000;
    total = (total + p65) % 10000000;
    total = (total + p66) % 10000000;
    total = (total + p67) % 10000000;
    total = (total + p68) % 10000000;
    total = (total + p69) % 10000000;
    total = (total + p70) % 10000000;
    total = (total + p71) % 10000000;
    total = (total + p72) % 10000000;
    total = (total + p73) % 10000000;
    total = (total + p74) % 10000000;
    total = (total + p75) % 10000000;
    total = (total + p76) % 10000000;
    total = (total + p77) % 10000000;
    total = (total + p78) % 10000000;
    total = (total + p79) % 10000000;
    total = (total + p80) % 10000000;
    total = (total + p81) % 10000000;
    total = (total + p82) % 10000000;
    total = (total + p83) % 10000000;
    total = (total + p84) % 10000000;
    total = (total + p85) % 10000000;
    total = (total + p86) % 10000000;
    total = (total + p87) % 10000000;
    total = (total + p88) % 10000000;
    total = (total + p89) % 10000000;
    total = (total + p90) % 10000000;
    total = (total + p91) % 10000000;
    total = (total + p92) % 10000000;
    total = (total + p93) % 10000000;
    total = (total + p94) % 10000000;
    total = (total + p95) % 10000000;
    total = (total + p96) % 10000000;
    total = (total + p97) % 10000000;
    total = (total + p98) % 10000000;
    total = (total + p99) % 10000000;
    return total;
}

function callSite0() returns int {
    return bigSum();
}
function callSite1() returns int {
    return bigSum(p7 = 8, p14 = 15, p34 = 35);
}
function callSite2() returns int {
    return bigSum(p14 = 16, p27 = 29, p63 = 65);
}
function callSite3() returns int {
    return bigSum(p21 = 24, p40 = 43, p92 = 95);
}
function callSite4() returns int {
    return bigSum(p28 = 32, p53 = 57, p21 = 25);
}
function callSite5() returns int {
    return bigSum(p35 = 40, p66 = 71, p50 = 55);
}
function callSite6() returns int {
    return bigSum(p42 = 48, p79 = 85);
}
function callSite7() returns int {
    return bigSum(p49 = 56, p92 = 99, p8 = 15);
}
function callSite8() returns int {
    return bigSum(p56 = 64, p5 = 13, p37 = 45);
}
function callSite9() returns int {
    return bigSum(p63 = 72, p18 = 27, p66 = 75);
}
function callSite10() returns int {
    return bigSum(p70 = 80, p31 = 41, p95 = 105);
}
function callSite11() returns int {
    return bigSum(p77 = 88, p44 = 55, p24 = 35);
}
function callSite12() returns int {
    return bigSum(p84 = 96, p57 = 69, p53 = 65);
}
function callSite13() returns int {
    return bigSum(p91 = 104, p70 = 83, p82 = 95);
}
function callSite14() returns int {
    return bigSum(p98 = 112, p83 = 97, p11 = 25);
}
function callSite15() returns int {
    return bigSum(p5 = 20, p96 = 111, p40 = 55);
}
function callSite16() returns int {
    return bigSum(p12 = 28, p9 = 25, p69 = 85);
}
function callSite17() returns int {
    return bigSum(p19 = 36, p22 = 39, p98 = 115);
}
function callSite18() returns int {
    return bigSum(p26 = 44, p35 = 53, p27 = 45);
}
function callSite19() returns int {
    return bigSum(p33 = 52, p48 = 67, p56 = 75);
}
function callSite20() returns int {
    return bigSum(p40 = 60, p61 = 81, p85 = 105);
}
function callSite21() returns int {
    return bigSum(p47 = 68, p74 = 95, p14 = 35);
}
function callSite22() returns int {
    return bigSum(p54 = 76, p87 = 109, p43 = 65);
}
function callSite23() returns int {
    return bigSum(p61 = 84, p0 = 23, p72 = 95);
}
function callSite24() returns int {
    return bigSum(p68 = 92, p13 = 37, p1 = 25);
}
function callSite25() returns int {
    return bigSum(p75 = 100, p26 = 51, p30 = 55);
}
function callSite26() returns int {
    return bigSum(p82 = 108, p39 = 65, p59 = 85);
}
function callSite27() returns int {
    return bigSum(p89 = 116, p52 = 79, p88 = 115);
}
function callSite28() returns int {
    return bigSum(p96 = 124, p65 = 93, p17 = 45);
}
function callSite29() returns int {
    return bigSum(p3 = 32, p78 = 107, p46 = 75);
}
function callSite30() returns int {
    return bigSum(p10 = 40, p91 = 121, p75 = 105);
}
function callSite31() returns int {
    return bigSum(p17 = 48, p4 = 35);
}
function callSite32() returns int {
    return bigSum(p24 = 56, p17 = 49, p33 = 65);
}
function callSite33() returns int {
    return bigSum(p31 = 64, p30 = 63, p62 = 95);
}
function callSite34() returns int {
    return bigSum(p38 = 72, p43 = 77, p91 = 125);
}
function callSite35() returns int {
    return bigSum(p45 = 80, p56 = 91, p20 = 55);
}
function callSite36() returns int {
    return bigSum(p52 = 88, p69 = 105, p49 = 85);
}
function callSite37() returns int {
    return bigSum(p59 = 96, p82 = 119, p78 = 115);
}
function callSite38() returns int {
    return bigSum(p66 = 104, p95 = 133, p7 = 45);
}
function callSite39() returns int {
    return bigSum(p73 = 112, p8 = 47, p36 = 75);
}
function callSite40() returns int {
    return bigSum(p80 = 120, p21 = 61, p65 = 105);
}
function callSite41() returns int {
    return bigSum(p87 = 128, p34 = 75, p94 = 135);
}
function callSite42() returns int {
    return bigSum(p94 = 136, p47 = 89, p23 = 65);
}
function callSite43() returns int {
    return bigSum(p1 = 44, p60 = 103, p52 = 95);
}
function callSite44() returns int {
    return bigSum(p8 = 52, p73 = 117, p81 = 125);
}
function callSite45() returns int {
    return bigSum(p15 = 60, p86 = 131, p10 = 55);
}
function callSite46() returns int {
    return bigSum(p22 = 68, p99 = 145, p39 = 85);
}
function callSite47() returns int {
    return bigSum(p29 = 76, p12 = 59, p68 = 115);
}
function callSite48() returns int {
    return bigSum(p36 = 84, p25 = 73, p97 = 145);
}
function callSite49() returns int {
    return bigSum(p43 = 92, p38 = 87, p26 = 75);
}

public function main() {
    int total = 0;
    total = (total + callSite0()) % 10000000;
    total = (total + callSite1()) % 10000000;
    total = (total + callSite2()) % 10000000;
    total = (total + callSite3()) % 10000000;
    total = (total + callSite4()) % 10000000;
    total = (total + callSite5()) % 10000000;
    total = (total + callSite6()) % 10000000;
    total = (total + callSite7()) % 10000000;
    total = (total + callSite8()) % 10000000;
    total = (total + callSite9()) % 10000000;
    total = (total + callSite10()) % 10000000;
    total = (total + callSite11()) % 10000000;
    total = (total + callSite12()) % 10000000;
    total = (total + callSite13()) % 10000000;
    total = (total + callSite14()) % 10000000;
    total = (total + callSite15()) % 10000000;
    total = (total + callSite16()) % 10000000;
    total = (total + callSite17()) % 10000000;
    total = (total + callSite18()) % 10000000;
    total = (total + callSite19()) % 10000000;
    total = (total + callSite20()) % 10000000;
    total = (total + callSite21()) % 10000000;
    total = (total + callSite22()) % 10000000;
    total = (total + callSite23()) % 10000000;
    total = (total + callSite24()) % 10000000;
    total = (total + callSite25()) % 10000000;
    total = (total + callSite26()) % 10000000;
    total = (total + callSite27()) % 10000000;
    total = (total + callSite28()) % 10000000;
    total = (total + callSite29()) % 10000000;
    total = (total + callSite30()) % 10000000;
    total = (total + callSite31()) % 10000000;
    total = (total + callSite32()) % 10000000;
    total = (total + callSite33()) % 10000000;
    total = (total + callSite34()) % 10000000;
    total = (total + callSite35()) % 10000000;
    total = (total + callSite36()) % 10000000;
    total = (total + callSite37()) % 10000000;
    total = (total + callSite38()) % 10000000;
    total = (total + callSite39()) % 10000000;
    total = (total + callSite40()) % 10000000;
    total = (total + callSite41()) % 10000000;
    total = (total + callSite42()) % 10000000;
    total = (total + callSite43()) % 10000000;
    total = (total + callSite44()) % 10000000;
    total = (total + callSite45()) % 10000000;
    total = (total + callSite46()) % 10000000;
    total = (total + callSite47()) % 10000000;
    total = (total + callSite48()) % 10000000;
    total = (total + callSite49()) % 10000000;
    io:println(total);
}
