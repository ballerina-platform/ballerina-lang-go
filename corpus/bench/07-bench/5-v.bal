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

const int C0 = 7;
const int C1 = 38;
const int C2 = 69;
const int C3 = 100;
const int C4 = 131;
const int C5 = 162;
const int C6 = 193;
const int C7 = 224;
const int C8 = 255;
const int C9 = 286;
const int C10 = 317;
const int C11 = 348;
const int C12 = 379;
const int C13 = 410;
const int C14 = 441;
const int C15 = 472;
const int C16 = 503;
const int C17 = 534;
const int C18 = 565;
const int C19 = 596;
const int C20 = 627;
const int C21 = 658;
const int C22 = 689;
const int C23 = 720;
const int C24 = 751;
const int C25 = 782;
const int C26 = 813;
const int C27 = 844;
const int C28 = 875;
const int C29 = 906;
const int C30 = 937;
const int C31 = 968;
const int C32 = 999;
const int C33 = 1030;
const int C34 = 1061;
const int C35 = 1092;
const int C36 = 1123;
const int C37 = 1154;
const int C38 = 1185;
const int C39 = 1216;
const int C40 = 1247;
const int C41 = 1278;
const int C42 = 1309;
const int C43 = 1340;
const int C44 = 1371;
const int C45 = 1402;
const int C46 = 1433;
const int C47 = 1464;
const int C48 = 1495;
const int C49 = 1526;
const int C50 = 1557;
const int C51 = 1588;
const int C52 = 1619;
const int C53 = 1650;
const int C54 = 1681;
const int C55 = 1712;
const int C56 = 1743;
const int C57 = 1774;
const int C58 = 1805;
const int C59 = 1836;
const int C60 = 1867;
const int C61 = 1898;
const int C62 = 1929;
const int C63 = 1960;

type R0 record {
    int id;
    int a = C0 + 0;
    int b = C1 * 1;
    int c = (0 + C0) % 10000;
    boolean flag = true;
    string label = "rec0";
};
type R1 record {
    int id;
    int a = C1 + 1;
    int b = C4 * 2;
    int c = (1 + C1) % 10000;
    boolean flag = false;
    string label = "rec1";
};
type R2 record {
    int id;
    int a = C2 + 2;
    int b = C7 * 3;
    int c = (2 + C2) % 10000;
    boolean flag = true;
    string label = "rec2";
};
type R3 record {
    int id;
    int a = C3 + 3;
    int b = C10 * 4;
    int c = (3 + C3) % 10000;
    boolean flag = false;
    string label = "rec3";
};
type R4 record {
    int id;
    int a = C4 + 4;
    int b = C13 * 5;
    int c = (4 + C4) % 10000;
    boolean flag = true;
    string label = "rec4";
};
type R5 record {
    int id;
    int a = C5 + 5;
    int b = C16 * 6;
    int c = (5 + C5) % 10000;
    boolean flag = false;
    string label = "rec5";
};
type R6 record {
    int id;
    int a = C6 + 6;
    int b = C19 * 7;
    int c = (6 + C6) % 10000;
    boolean flag = true;
    string label = "rec6";
};
type R7 record {
    int id;
    int a = C7 + 7;
    int b = C22 * 1;
    int c = (7 + C7) % 10000;
    boolean flag = false;
    string label = "rec7";
};
type R8 record {
    int id;
    int a = C8 + 8;
    int b = C25 * 2;
    int c = (8 + C8) % 10000;
    boolean flag = true;
    string label = "rec8";
};
type R9 record {
    int id;
    int a = C9 + 9;
    int b = C28 * 3;
    int c = (9 + C9) % 10000;
    boolean flag = false;
    string label = "rec9";
};
type R10 record {
    int id;
    int a = C10 + 10;
    int b = C31 * 4;
    int c = (10 + C10) % 10000;
    boolean flag = true;
    string label = "rec10";
};
type R11 record {
    int id;
    int a = C11 + 11;
    int b = C34 * 5;
    int c = (11 + C11) % 10000;
    boolean flag = false;
    string label = "rec11";
};
type R12 record {
    int id;
    int a = C12 + 12;
    int b = C37 * 6;
    int c = (12 + C12) % 10000;
    boolean flag = true;
    string label = "rec12";
};
type R13 record {
    int id;
    int a = C13 + 13;
    int b = C40 * 7;
    int c = (13 + C13) % 10000;
    boolean flag = false;
    string label = "rec13";
};
type R14 record {
    int id;
    int a = C14 + 14;
    int b = C43 * 1;
    int c = (14 + C14) % 10000;
    boolean flag = true;
    string label = "rec14";
};
type R15 record {
    int id;
    int a = C15 + 15;
    int b = C46 * 2;
    int c = (15 + C15) % 10000;
    boolean flag = false;
    string label = "rec15";
};
type R16 record {
    int id;
    int a = C16 + 16;
    int b = C49 * 3;
    int c = (16 + C16) % 10000;
    boolean flag = true;
    string label = "rec16";
};
type R17 record {
    int id;
    int a = C17 + 17;
    int b = C52 * 4;
    int c = (17 + C17) % 10000;
    boolean flag = false;
    string label = "rec17";
};
type R18 record {
    int id;
    int a = C18 + 18;
    int b = C55 * 5;
    int c = (18 + C18) % 10000;
    boolean flag = true;
    string label = "rec18";
};
type R19 record {
    int id;
    int a = C19 + 19;
    int b = C58 * 6;
    int c = (19 + C19) % 10000;
    boolean flag = false;
    string label = "rec19";
};
type R20 record {
    int id;
    int a = C20 + 20;
    int b = C61 * 7;
    int c = (20 + C20) % 10000;
    boolean flag = true;
    string label = "rec20";
};
type R21 record {
    int id;
    int a = C21 + 21;
    int b = C0 * 1;
    int c = (21 + C21) % 10000;
    boolean flag = false;
    string label = "rec21";
};
type R22 record {
    int id;
    int a = C22 + 22;
    int b = C3 * 2;
    int c = (22 + C22) % 10000;
    boolean flag = true;
    string label = "rec22";
};
type R23 record {
    int id;
    int a = C23 + 23;
    int b = C6 * 3;
    int c = (23 + C23) % 10000;
    boolean flag = false;
    string label = "rec23";
};
type R24 record {
    int id;
    int a = C24 + 24;
    int b = C9 * 4;
    int c = (24 + C24) % 10000;
    boolean flag = true;
    string label = "rec24";
};
type R25 record {
    int id;
    int a = C25 + 25;
    int b = C12 * 5;
    int c = (25 + C25) % 10000;
    boolean flag = false;
    string label = "rec25";
};
type R26 record {
    int id;
    int a = C26 + 26;
    int b = C15 * 6;
    int c = (26 + C26) % 10000;
    boolean flag = true;
    string label = "rec26";
};
type R27 record {
    int id;
    int a = C27 + 27;
    int b = C18 * 7;
    int c = (27 + C27) % 10000;
    boolean flag = false;
    string label = "rec27";
};
type R28 record {
    int id;
    int a = C28 + 28;
    int b = C21 * 1;
    int c = (28 + C28) % 10000;
    boolean flag = true;
    string label = "rec28";
};
type R29 record {
    int id;
    int a = C29 + 29;
    int b = C24 * 2;
    int c = (29 + C29) % 10000;
    boolean flag = false;
    string label = "rec29";
};
type R30 record {
    int id;
    int a = C30 + 30;
    int b = C27 * 3;
    int c = (30 + C30) % 10000;
    boolean flag = true;
    string label = "rec30";
};
type R31 record {
    int id;
    int a = C31 + 31;
    int b = C30 * 4;
    int c = (31 + C31) % 10000;
    boolean flag = false;
    string label = "rec31";
};
type R32 record {
    int id;
    int a = C32 + 32;
    int b = C33 * 5;
    int c = (32 + C32) % 10000;
    boolean flag = true;
    string label = "rec32";
};
type R33 record {
    int id;
    int a = C33 + 33;
    int b = C36 * 6;
    int c = (33 + C33) % 10000;
    boolean flag = false;
    string label = "rec33";
};
type R34 record {
    int id;
    int a = C34 + 34;
    int b = C39 * 7;
    int c = (34 + C34) % 10000;
    boolean flag = true;
    string label = "rec34";
};
type R35 record {
    int id;
    int a = C35 + 35;
    int b = C42 * 1;
    int c = (35 + C35) % 10000;
    boolean flag = false;
    string label = "rec35";
};
type R36 record {
    int id;
    int a = C36 + 36;
    int b = C45 * 2;
    int c = (36 + C36) % 10000;
    boolean flag = true;
    string label = "rec36";
};
type R37 record {
    int id;
    int a = C37 + 37;
    int b = C48 * 3;
    int c = (37 + C37) % 10000;
    boolean flag = false;
    string label = "rec37";
};
type R38 record {
    int id;
    int a = C38 + 38;
    int b = C51 * 4;
    int c = (38 + C38) % 10000;
    boolean flag = true;
    string label = "rec38";
};
type R39 record {
    int id;
    int a = C39 + 39;
    int b = C54 * 5;
    int c = (39 + C39) % 10000;
    boolean flag = false;
    string label = "rec39";
};
type R40 record {
    int id;
    int a = C40 + 40;
    int b = C57 * 6;
    int c = (40 + C40) % 10000;
    boolean flag = true;
    string label = "rec40";
};
type R41 record {
    int id;
    int a = C41 + 41;
    int b = C60 * 7;
    int c = (41 + C41) % 10000;
    boolean flag = false;
    string label = "rec41";
};
type R42 record {
    int id;
    int a = C42 + 42;
    int b = C63 * 1;
    int c = (42 + C42) % 10000;
    boolean flag = true;
    string label = "rec42";
};
type R43 record {
    int id;
    int a = C43 + 43;
    int b = C2 * 2;
    int c = (43 + C43) % 10000;
    boolean flag = false;
    string label = "rec43";
};
type R44 record {
    int id;
    int a = C44 + 44;
    int b = C5 * 3;
    int c = (44 + C44) % 10000;
    boolean flag = true;
    string label = "rec44";
};
type R45 record {
    int id;
    int a = C45 + 45;
    int b = C8 * 4;
    int c = (45 + C45) % 10000;
    boolean flag = false;
    string label = "rec45";
};
type R46 record {
    int id;
    int a = C46 + 46;
    int b = C11 * 5;
    int c = (46 + C46) % 10000;
    boolean flag = true;
    string label = "rec46";
};
type R47 record {
    int id;
    int a = C47 + 47;
    int b = C14 * 6;
    int c = (47 + C47) % 10000;
    boolean flag = false;
    string label = "rec47";
};
type R48 record {
    int id;
    int a = C48 + 48;
    int b = C17 * 7;
    int c = (48 + C48) % 10000;
    boolean flag = true;
    string label = "rec48";
};
type R49 record {
    int id;
    int a = C49 + 49;
    int b = C20 * 1;
    int c = (49 + C49) % 10000;
    boolean flag = false;
    string label = "rec49";
};
type R50 record {
    int id;
    int a = C50 + 50;
    int b = C23 * 2;
    int c = (50 + C50) % 10000;
    boolean flag = true;
    string label = "rec50";
};
type R51 record {
    int id;
    int a = C51 + 51;
    int b = C26 * 3;
    int c = (51 + C51) % 10000;
    boolean flag = false;
    string label = "rec51";
};
type R52 record {
    int id;
    int a = C52 + 52;
    int b = C29 * 4;
    int c = (52 + C52) % 10000;
    boolean flag = true;
    string label = "rec52";
};
type R53 record {
    int id;
    int a = C53 + 53;
    int b = C32 * 5;
    int c = (53 + C53) % 10000;
    boolean flag = false;
    string label = "rec53";
};
type R54 record {
    int id;
    int a = C54 + 54;
    int b = C35 * 6;
    int c = (54 + C54) % 10000;
    boolean flag = true;
    string label = "rec54";
};
type R55 record {
    int id;
    int a = C55 + 55;
    int b = C38 * 7;
    int c = (55 + C55) % 10000;
    boolean flag = false;
    string label = "rec55";
};
type R56 record {
    int id;
    int a = C56 + 56;
    int b = C41 * 1;
    int c = (56 + C56) % 10000;
    boolean flag = true;
    string label = "rec56";
};
type R57 record {
    int id;
    int a = C57 + 57;
    int b = C44 * 2;
    int c = (57 + C57) % 10000;
    boolean flag = false;
    string label = "rec57";
};
type R58 record {
    int id;
    int a = C58 + 58;
    int b = C47 * 3;
    int c = (58 + C58) % 10000;
    boolean flag = true;
    string label = "rec58";
};
type R59 record {
    int id;
    int a = C59 + 59;
    int b = C50 * 4;
    int c = (59 + C59) % 10000;
    boolean flag = false;
    string label = "rec59";
};
type R60 record {
    int id;
    int a = C60 + 60;
    int b = C53 * 5;
    int c = (60 + C60) % 10000;
    boolean flag = true;
    string label = "rec60";
};
type R61 record {
    int id;
    int a = C61 + 61;
    int b = C56 * 6;
    int c = (61 + C61) % 10000;
    boolean flag = false;
    string label = "rec61";
};
type R62 record {
    int id;
    int a = C62 + 62;
    int b = C59 * 7;
    int c = (62 + C62) % 10000;
    boolean flag = true;
    string label = "rec62";
};
type R63 record {
    int id;
    int a = C63 + 63;
    int b = C62 * 1;
    int c = (63 + C63) % 10000;
    boolean flag = false;
    string label = "rec63";
};
type R64 record {
    int id;
    int a = C0 + 64;
    int b = C1 * 2;
    int c = (64 + C0) % 10000;
    boolean flag = true;
    string label = "rec64";
};
type R65 record {
    int id;
    int a = C1 + 65;
    int b = C4 * 3;
    int c = (65 + C1) % 10000;
    boolean flag = false;
    string label = "rec65";
};
type R66 record {
    int id;
    int a = C2 + 66;
    int b = C7 * 4;
    int c = (66 + C2) % 10000;
    boolean flag = true;
    string label = "rec66";
};
type R67 record {
    int id;
    int a = C3 + 67;
    int b = C10 * 5;
    int c = (67 + C3) % 10000;
    boolean flag = false;
    string label = "rec67";
};
type R68 record {
    int id;
    int a = C4 + 68;
    int b = C13 * 6;
    int c = (68 + C4) % 10000;
    boolean flag = true;
    string label = "rec68";
};
type R69 record {
    int id;
    int a = C5 + 69;
    int b = C16 * 7;
    int c = (69 + C5) % 10000;
    boolean flag = false;
    string label = "rec69";
};
type R70 record {
    int id;
    int a = C6 + 70;
    int b = C19 * 1;
    int c = (70 + C6) % 10000;
    boolean flag = true;
    string label = "rec70";
};
type R71 record {
    int id;
    int a = C7 + 71;
    int b = C22 * 2;
    int c = (71 + C7) % 10000;
    boolean flag = false;
    string label = "rec71";
};
type R72 record {
    int id;
    int a = C8 + 72;
    int b = C25 * 3;
    int c = (72 + C8) % 10000;
    boolean flag = true;
    string label = "rec72";
};
type R73 record {
    int id;
    int a = C9 + 73;
    int b = C28 * 4;
    int c = (73 + C9) % 10000;
    boolean flag = false;
    string label = "rec73";
};
type R74 record {
    int id;
    int a = C10 + 74;
    int b = C31 * 5;
    int c = (74 + C10) % 10000;
    boolean flag = true;
    string label = "rec74";
};
type R75 record {
    int id;
    int a = C11 + 75;
    int b = C34 * 6;
    int c = (75 + C11) % 10000;
    boolean flag = false;
    string label = "rec75";
};
type R76 record {
    int id;
    int a = C12 + 76;
    int b = C37 * 7;
    int c = (76 + C12) % 10000;
    boolean flag = true;
    string label = "rec76";
};
type R77 record {
    int id;
    int a = C13 + 77;
    int b = C40 * 1;
    int c = (77 + C13) % 10000;
    boolean flag = false;
    string label = "rec77";
};
type R78 record {
    int id;
    int a = C14 + 78;
    int b = C43 * 2;
    int c = (78 + C14) % 10000;
    boolean flag = true;
    string label = "rec78";
};
type R79 record {
    int id;
    int a = C15 + 79;
    int b = C46 * 3;
    int c = (79 + C15) % 10000;
    boolean flag = false;
    string label = "rec79";
};
type R80 record {
    int id;
    int a = C16 + 80;
    int b = C49 * 4;
    int c = (80 + C16) % 10000;
    boolean flag = true;
    string label = "rec80";
};
type R81 record {
    int id;
    int a = C17 + 81;
    int b = C52 * 5;
    int c = (81 + C17) % 10000;
    boolean flag = false;
    string label = "rec81";
};
type R82 record {
    int id;
    int a = C18 + 82;
    int b = C55 * 6;
    int c = (82 + C18) % 10000;
    boolean flag = true;
    string label = "rec82";
};
type R83 record {
    int id;
    int a = C19 + 83;
    int b = C58 * 7;
    int c = (83 + C19) % 10000;
    boolean flag = false;
    string label = "rec83";
};
type R84 record {
    int id;
    int a = C20 + 84;
    int b = C61 * 1;
    int c = (84 + C20) % 10000;
    boolean flag = true;
    string label = "rec84";
};
type R85 record {
    int id;
    int a = C21 + 85;
    int b = C0 * 2;
    int c = (85 + C21) % 10000;
    boolean flag = false;
    string label = "rec85";
};
type R86 record {
    int id;
    int a = C22 + 86;
    int b = C3 * 3;
    int c = (86 + C22) % 10000;
    boolean flag = true;
    string label = "rec86";
};
type R87 record {
    int id;
    int a = C23 + 87;
    int b = C6 * 4;
    int c = (87 + C23) % 10000;
    boolean flag = false;
    string label = "rec87";
};
type R88 record {
    int id;
    int a = C24 + 88;
    int b = C9 * 5;
    int c = (88 + C24) % 10000;
    boolean flag = true;
    string label = "rec88";
};
type R89 record {
    int id;
    int a = C25 + 89;
    int b = C12 * 6;
    int c = (89 + C25) % 10000;
    boolean flag = false;
    string label = "rec89";
};
type R90 record {
    int id;
    int a = C26 + 90;
    int b = C15 * 7;
    int c = (90 + C26) % 10000;
    boolean flag = true;
    string label = "rec90";
};
type R91 record {
    int id;
    int a = C27 + 91;
    int b = C18 * 1;
    int c = (91 + C27) % 10000;
    boolean flag = false;
    string label = "rec91";
};
type R92 record {
    int id;
    int a = C28 + 92;
    int b = C21 * 2;
    int c = (92 + C28) % 10000;
    boolean flag = true;
    string label = "rec92";
};
type R93 record {
    int id;
    int a = C29 + 93;
    int b = C24 * 3;
    int c = (93 + C29) % 10000;
    boolean flag = false;
    string label = "rec93";
};
type R94 record {
    int id;
    int a = C30 + 94;
    int b = C27 * 4;
    int c = (94 + C30) % 10000;
    boolean flag = true;
    string label = "rec94";
};
type R95 record {
    int id;
    int a = C31 + 95;
    int b = C30 * 5;
    int c = (95 + C31) % 10000;
    boolean flag = false;
    string label = "rec95";
};
type R96 record {
    int id;
    int a = C32 + 96;
    int b = C33 * 6;
    int c = (96 + C32) % 10000;
    boolean flag = true;
    string label = "rec96";
};
type R97 record {
    int id;
    int a = C33 + 97;
    int b = C36 * 7;
    int c = (97 + C33) % 10000;
    boolean flag = false;
    string label = "rec97";
};
type R98 record {
    int id;
    int a = C34 + 98;
    int b = C39 * 1;
    int c = (98 + C34) % 10000;
    boolean flag = true;
    string label = "rec98";
};
type R99 record {
    int id;
    int a = C35 + 99;
    int b = C42 * 2;
    int c = (99 + C35) % 10000;
    boolean flag = false;
    string label = "rec99";
};
type R100 record {
    int id;
    int a = C36 + 100;
    int b = C45 * 3;
    int c = (100 + C36) % 10000;
    boolean flag = true;
    string label = "rec100";
};
type R101 record {
    int id;
    int a = C37 + 101;
    int b = C48 * 4;
    int c = (101 + C37) % 10000;
    boolean flag = false;
    string label = "rec101";
};
type R102 record {
    int id;
    int a = C38 + 102;
    int b = C51 * 5;
    int c = (102 + C38) % 10000;
    boolean flag = true;
    string label = "rec102";
};
type R103 record {
    int id;
    int a = C39 + 103;
    int b = C54 * 6;
    int c = (103 + C39) % 10000;
    boolean flag = false;
    string label = "rec103";
};
type R104 record {
    int id;
    int a = C40 + 104;
    int b = C57 * 7;
    int c = (104 + C40) % 10000;
    boolean flag = true;
    string label = "rec104";
};
type R105 record {
    int id;
    int a = C41 + 105;
    int b = C60 * 1;
    int c = (105 + C41) % 10000;
    boolean flag = false;
    string label = "rec105";
};
type R106 record {
    int id;
    int a = C42 + 106;
    int b = C63 * 2;
    int c = (106 + C42) % 10000;
    boolean flag = true;
    string label = "rec106";
};
type R107 record {
    int id;
    int a = C43 + 107;
    int b = C2 * 3;
    int c = (107 + C43) % 10000;
    boolean flag = false;
    string label = "rec107";
};
type R108 record {
    int id;
    int a = C44 + 108;
    int b = C5 * 4;
    int c = (108 + C44) % 10000;
    boolean flag = true;
    string label = "rec108";
};
type R109 record {
    int id;
    int a = C45 + 109;
    int b = C8 * 5;
    int c = (109 + C45) % 10000;
    boolean flag = false;
    string label = "rec109";
};
type R110 record {
    int id;
    int a = C46 + 110;
    int b = C11 * 6;
    int c = (110 + C46) % 10000;
    boolean flag = true;
    string label = "rec110";
};
type R111 record {
    int id;
    int a = C47 + 111;
    int b = C14 * 7;
    int c = (111 + C47) % 10000;
    boolean flag = false;
    string label = "rec111";
};
type R112 record {
    int id;
    int a = C48 + 112;
    int b = C17 * 1;
    int c = (112 + C48) % 10000;
    boolean flag = true;
    string label = "rec112";
};
type R113 record {
    int id;
    int a = C49 + 113;
    int b = C20 * 2;
    int c = (113 + C49) % 10000;
    boolean flag = false;
    string label = "rec113";
};
type R114 record {
    int id;
    int a = C50 + 114;
    int b = C23 * 3;
    int c = (114 + C50) % 10000;
    boolean flag = true;
    string label = "rec114";
};
type R115 record {
    int id;
    int a = C51 + 115;
    int b = C26 * 4;
    int c = (115 + C51) % 10000;
    boolean flag = false;
    string label = "rec115";
};
type R116 record {
    int id;
    int a = C52 + 116;
    int b = C29 * 5;
    int c = (116 + C52) % 10000;
    boolean flag = true;
    string label = "rec116";
};
type R117 record {
    int id;
    int a = C53 + 117;
    int b = C32 * 6;
    int c = (117 + C53) % 10000;
    boolean flag = false;
    string label = "rec117";
};
type R118 record {
    int id;
    int a = C54 + 118;
    int b = C35 * 7;
    int c = (118 + C54) % 10000;
    boolean flag = true;
    string label = "rec118";
};
type R119 record {
    int id;
    int a = C55 + 119;
    int b = C38 * 1;
    int c = (119 + C55) % 10000;
    boolean flag = false;
    string label = "rec119";
};
type R120 record {
    int id;
    int a = C56 + 120;
    int b = C41 * 2;
    int c = (120 + C56) % 10000;
    boolean flag = true;
    string label = "rec120";
};
type R121 record {
    int id;
    int a = C57 + 121;
    int b = C44 * 3;
    int c = (121 + C57) % 10000;
    boolean flag = false;
    string label = "rec121";
};
type R122 record {
    int id;
    int a = C58 + 122;
    int b = C47 * 4;
    int c = (122 + C58) % 10000;
    boolean flag = true;
    string label = "rec122";
};
type R123 record {
    int id;
    int a = C59 + 123;
    int b = C50 * 5;
    int c = (123 + C59) % 10000;
    boolean flag = false;
    string label = "rec123";
};
type R124 record {
    int id;
    int a = C60 + 124;
    int b = C53 * 6;
    int c = (124 + C60) % 10000;
    boolean flag = true;
    string label = "rec124";
};
type R125 record {
    int id;
    int a = C61 + 125;
    int b = C56 * 7;
    int c = (125 + C61) % 10000;
    boolean flag = false;
    string label = "rec125";
};
type R126 record {
    int id;
    int a = C62 + 126;
    int b = C59 * 1;
    int c = (126 + C62) % 10000;
    boolean flag = true;
    string label = "rec126";
};
type R127 record {
    int id;
    int a = C63 + 127;
    int b = C62 * 2;
    int c = (127 + C63) % 10000;
    boolean flag = false;
    string label = "rec127";
};
type R128 record {
    int id;
    int a = C0 + 128;
    int b = C1 * 3;
    int c = (128 + C0) % 10000;
    boolean flag = true;
    string label = "rec128";
};
type R129 record {
    int id;
    int a = C1 + 129;
    int b = C4 * 4;
    int c = (129 + C1) % 10000;
    boolean flag = false;
    string label = "rec129";
};
type R130 record {
    int id;
    int a = C2 + 130;
    int b = C7 * 5;
    int c = (130 + C2) % 10000;
    boolean flag = true;
    string label = "rec130";
};
type R131 record {
    int id;
    int a = C3 + 131;
    int b = C10 * 6;
    int c = (131 + C3) % 10000;
    boolean flag = false;
    string label = "rec131";
};
type R132 record {
    int id;
    int a = C4 + 132;
    int b = C13 * 7;
    int c = (132 + C4) % 10000;
    boolean flag = true;
    string label = "rec132";
};
type R133 record {
    int id;
    int a = C5 + 133;
    int b = C16 * 1;
    int c = (133 + C5) % 10000;
    boolean flag = false;
    string label = "rec133";
};
type R134 record {
    int id;
    int a = C6 + 134;
    int b = C19 * 2;
    int c = (134 + C6) % 10000;
    boolean flag = true;
    string label = "rec134";
};
type R135 record {
    int id;
    int a = C7 + 135;
    int b = C22 * 3;
    int c = (135 + C7) % 10000;
    boolean flag = false;
    string label = "rec135";
};
type R136 record {
    int id;
    int a = C8 + 136;
    int b = C25 * 4;
    int c = (136 + C8) % 10000;
    boolean flag = true;
    string label = "rec136";
};
type R137 record {
    int id;
    int a = C9 + 137;
    int b = C28 * 5;
    int c = (137 + C9) % 10000;
    boolean flag = false;
    string label = "rec137";
};
type R138 record {
    int id;
    int a = C10 + 138;
    int b = C31 * 6;
    int c = (138 + C10) % 10000;
    boolean flag = true;
    string label = "rec138";
};
type R139 record {
    int id;
    int a = C11 + 139;
    int b = C34 * 7;
    int c = (139 + C11) % 10000;
    boolean flag = false;
    string label = "rec139";
};
type R140 record {
    int id;
    int a = C12 + 140;
    int b = C37 * 1;
    int c = (140 + C12) % 10000;
    boolean flag = true;
    string label = "rec140";
};
type R141 record {
    int id;
    int a = C13 + 141;
    int b = C40 * 2;
    int c = (141 + C13) % 10000;
    boolean flag = false;
    string label = "rec141";
};
type R142 record {
    int id;
    int a = C14 + 142;
    int b = C43 * 3;
    int c = (142 + C14) % 10000;
    boolean flag = true;
    string label = "rec142";
};
type R143 record {
    int id;
    int a = C15 + 143;
    int b = C46 * 4;
    int c = (143 + C15) % 10000;
    boolean flag = false;
    string label = "rec143";
};
type R144 record {
    int id;
    int a = C16 + 144;
    int b = C49 * 5;
    int c = (144 + C16) % 10000;
    boolean flag = true;
    string label = "rec144";
};
type R145 record {
    int id;
    int a = C17 + 145;
    int b = C52 * 6;
    int c = (145 + C17) % 10000;
    boolean flag = false;
    string label = "rec145";
};
type R146 record {
    int id;
    int a = C18 + 146;
    int b = C55 * 7;
    int c = (146 + C18) % 10000;
    boolean flag = true;
    string label = "rec146";
};
type R147 record {
    int id;
    int a = C19 + 147;
    int b = C58 * 1;
    int c = (147 + C19) % 10000;
    boolean flag = false;
    string label = "rec147";
};
type R148 record {
    int id;
    int a = C20 + 148;
    int b = C61 * 2;
    int c = (148 + C20) % 10000;
    boolean flag = true;
    string label = "rec148";
};
type R149 record {
    int id;
    int a = C21 + 149;
    int b = C0 * 3;
    int c = (149 + C21) % 10000;
    boolean flag = false;
    string label = "rec149";
};
type R150 record {
    int id;
    int a = C22 + 150;
    int b = C3 * 4;
    int c = (150 + C22) % 10000;
    boolean flag = true;
    string label = "rec150";
};
type R151 record {
    int id;
    int a = C23 + 151;
    int b = C6 * 5;
    int c = (151 + C23) % 10000;
    boolean flag = false;
    string label = "rec151";
};
type R152 record {
    int id;
    int a = C24 + 152;
    int b = C9 * 6;
    int c = (152 + C24) % 10000;
    boolean flag = true;
    string label = "rec152";
};
type R153 record {
    int id;
    int a = C25 + 153;
    int b = C12 * 7;
    int c = (153 + C25) % 10000;
    boolean flag = false;
    string label = "rec153";
};
type R154 record {
    int id;
    int a = C26 + 154;
    int b = C15 * 1;
    int c = (154 + C26) % 10000;
    boolean flag = true;
    string label = "rec154";
};
type R155 record {
    int id;
    int a = C27 + 155;
    int b = C18 * 2;
    int c = (155 + C27) % 10000;
    boolean flag = false;
    string label = "rec155";
};
type R156 record {
    int id;
    int a = C28 + 156;
    int b = C21 * 3;
    int c = (156 + C28) % 10000;
    boolean flag = true;
    string label = "rec156";
};
type R157 record {
    int id;
    int a = C29 + 157;
    int b = C24 * 4;
    int c = (157 + C29) % 10000;
    boolean flag = false;
    string label = "rec157";
};
type R158 record {
    int id;
    int a = C30 + 158;
    int b = C27 * 5;
    int c = (158 + C30) % 10000;
    boolean flag = true;
    string label = "rec158";
};
type R159 record {
    int id;
    int a = C31 + 159;
    int b = C30 * 6;
    int c = (159 + C31) % 10000;
    boolean flag = false;
    string label = "rec159";
};
type R160 record {
    int id;
    int a = C32 + 160;
    int b = C33 * 7;
    int c = (160 + C32) % 10000;
    boolean flag = true;
    string label = "rec160";
};
type R161 record {
    int id;
    int a = C33 + 161;
    int b = C36 * 1;
    int c = (161 + C33) % 10000;
    boolean flag = false;
    string label = "rec161";
};
type R162 record {
    int id;
    int a = C34 + 162;
    int b = C39 * 2;
    int c = (162 + C34) % 10000;
    boolean flag = true;
    string label = "rec162";
};
type R163 record {
    int id;
    int a = C35 + 163;
    int b = C42 * 3;
    int c = (163 + C35) % 10000;
    boolean flag = false;
    string label = "rec163";
};
type R164 record {
    int id;
    int a = C36 + 164;
    int b = C45 * 4;
    int c = (164 + C36) % 10000;
    boolean flag = true;
    string label = "rec164";
};
type R165 record {
    int id;
    int a = C37 + 165;
    int b = C48 * 5;
    int c = (165 + C37) % 10000;
    boolean flag = false;
    string label = "rec165";
};
type R166 record {
    int id;
    int a = C38 + 166;
    int b = C51 * 6;
    int c = (166 + C38) % 10000;
    boolean flag = true;
    string label = "rec166";
};
type R167 record {
    int id;
    int a = C39 + 167;
    int b = C54 * 7;
    int c = (167 + C39) % 10000;
    boolean flag = false;
    string label = "rec167";
};
type R168 record {
    int id;
    int a = C40 + 168;
    int b = C57 * 1;
    int c = (168 + C40) % 10000;
    boolean flag = true;
    string label = "rec168";
};
type R169 record {
    int id;
    int a = C41 + 169;
    int b = C60 * 2;
    int c = (169 + C41) % 10000;
    boolean flag = false;
    string label = "rec169";
};
type R170 record {
    int id;
    int a = C42 + 170;
    int b = C63 * 3;
    int c = (170 + C42) % 10000;
    boolean flag = true;
    string label = "rec170";
};
type R171 record {
    int id;
    int a = C43 + 171;
    int b = C2 * 4;
    int c = (171 + C43) % 10000;
    boolean flag = false;
    string label = "rec171";
};
type R172 record {
    int id;
    int a = C44 + 172;
    int b = C5 * 5;
    int c = (172 + C44) % 10000;
    boolean flag = true;
    string label = "rec172";
};
type R173 record {
    int id;
    int a = C45 + 173;
    int b = C8 * 6;
    int c = (173 + C45) % 10000;
    boolean flag = false;
    string label = "rec173";
};
type R174 record {
    int id;
    int a = C46 + 174;
    int b = C11 * 7;
    int c = (174 + C46) % 10000;
    boolean flag = true;
    string label = "rec174";
};
type R175 record {
    int id;
    int a = C47 + 175;
    int b = C14 * 1;
    int c = (175 + C47) % 10000;
    boolean flag = false;
    string label = "rec175";
};
type R176 record {
    int id;
    int a = C48 + 176;
    int b = C17 * 2;
    int c = (176 + C48) % 10000;
    boolean flag = true;
    string label = "rec176";
};
type R177 record {
    int id;
    int a = C49 + 177;
    int b = C20 * 3;
    int c = (177 + C49) % 10000;
    boolean flag = false;
    string label = "rec177";
};
type R178 record {
    int id;
    int a = C50 + 178;
    int b = C23 * 4;
    int c = (178 + C50) % 10000;
    boolean flag = true;
    string label = "rec178";
};
type R179 record {
    int id;
    int a = C51 + 179;
    int b = C26 * 5;
    int c = (179 + C51) % 10000;
    boolean flag = false;
    string label = "rec179";
};
type R180 record {
    int id;
    int a = C52 + 180;
    int b = C29 * 6;
    int c = (180 + C52) % 10000;
    boolean flag = true;
    string label = "rec180";
};
type R181 record {
    int id;
    int a = C53 + 181;
    int b = C32 * 7;
    int c = (181 + C53) % 10000;
    boolean flag = false;
    string label = "rec181";
};
type R182 record {
    int id;
    int a = C54 + 182;
    int b = C35 * 1;
    int c = (182 + C54) % 10000;
    boolean flag = true;
    string label = "rec182";
};
type R183 record {
    int id;
    int a = C55 + 183;
    int b = C38 * 2;
    int c = (183 + C55) % 10000;
    boolean flag = false;
    string label = "rec183";
};
type R184 record {
    int id;
    int a = C56 + 184;
    int b = C41 * 3;
    int c = (184 + C56) % 10000;
    boolean flag = true;
    string label = "rec184";
};
type R185 record {
    int id;
    int a = C57 + 185;
    int b = C44 * 4;
    int c = (185 + C57) % 10000;
    boolean flag = false;
    string label = "rec185";
};
type R186 record {
    int id;
    int a = C58 + 186;
    int b = C47 * 5;
    int c = (186 + C58) % 10000;
    boolean flag = true;
    string label = "rec186";
};
type R187 record {
    int id;
    int a = C59 + 187;
    int b = C50 * 6;
    int c = (187 + C59) % 10000;
    boolean flag = false;
    string label = "rec187";
};
type R188 record {
    int id;
    int a = C60 + 188;
    int b = C53 * 7;
    int c = (188 + C60) % 10000;
    boolean flag = true;
    string label = "rec188";
};
type R189 record {
    int id;
    int a = C61 + 189;
    int b = C56 * 1;
    int c = (189 + C61) % 10000;
    boolean flag = false;
    string label = "rec189";
};
type R190 record {
    int id;
    int a = C62 + 190;
    int b = C59 * 2;
    int c = (190 + C62) % 10000;
    boolean flag = true;
    string label = "rec190";
};
type R191 record {
    int id;
    int a = C63 + 191;
    int b = C62 * 3;
    int c = (191 + C63) % 10000;
    boolean flag = false;
    string label = "rec191";
};
type R192 record {
    int id;
    int a = C0 + 192;
    int b = C1 * 4;
    int c = (192 + C0) % 10000;
    boolean flag = true;
    string label = "rec192";
};
type R193 record {
    int id;
    int a = C1 + 193;
    int b = C4 * 5;
    int c = (193 + C1) % 10000;
    boolean flag = false;
    string label = "rec193";
};
type R194 record {
    int id;
    int a = C2 + 194;
    int b = C7 * 6;
    int c = (194 + C2) % 10000;
    boolean flag = true;
    string label = "rec194";
};
type R195 record {
    int id;
    int a = C3 + 195;
    int b = C10 * 7;
    int c = (195 + C3) % 10000;
    boolean flag = false;
    string label = "rec195";
};
type R196 record {
    int id;
    int a = C4 + 196;
    int b = C13 * 1;
    int c = (196 + C4) % 10000;
    boolean flag = true;
    string label = "rec196";
};
type R197 record {
    int id;
    int a = C5 + 197;
    int b = C16 * 2;
    int c = (197 + C5) % 10000;
    boolean flag = false;
    string label = "rec197";
};
type R198 record {
    int id;
    int a = C6 + 198;
    int b = C19 * 3;
    int c = (198 + C6) % 10000;
    boolean flag = true;
    string label = "rec198";
};
type R199 record {
    int id;
    int a = C7 + 199;
    int b = C22 * 4;
    int c = (199 + C7) % 10000;
    boolean flag = false;
    string label = "rec199";
};
type R200 record {
    int id;
    int a = C8 + 200;
    int b = C25 * 5;
    int c = (200 + C8) % 10000;
    boolean flag = true;
    string label = "rec200";
};
type R201 record {
    int id;
    int a = C9 + 201;
    int b = C28 * 6;
    int c = (201 + C9) % 10000;
    boolean flag = false;
    string label = "rec201";
};
type R202 record {
    int id;
    int a = C10 + 202;
    int b = C31 * 7;
    int c = (202 + C10) % 10000;
    boolean flag = true;
    string label = "rec202";
};
type R203 record {
    int id;
    int a = C11 + 203;
    int b = C34 * 1;
    int c = (203 + C11) % 10000;
    boolean flag = false;
    string label = "rec203";
};
type R204 record {
    int id;
    int a = C12 + 204;
    int b = C37 * 2;
    int c = (204 + C12) % 10000;
    boolean flag = true;
    string label = "rec204";
};
type R205 record {
    int id;
    int a = C13 + 205;
    int b = C40 * 3;
    int c = (205 + C13) % 10000;
    boolean flag = false;
    string label = "rec205";
};
type R206 record {
    int id;
    int a = C14 + 206;
    int b = C43 * 4;
    int c = (206 + C14) % 10000;
    boolean flag = true;
    string label = "rec206";
};
type R207 record {
    int id;
    int a = C15 + 207;
    int b = C46 * 5;
    int c = (207 + C15) % 10000;
    boolean flag = false;
    string label = "rec207";
};
type R208 record {
    int id;
    int a = C16 + 208;
    int b = C49 * 6;
    int c = (208 + C16) % 10000;
    boolean flag = true;
    string label = "rec208";
};
type R209 record {
    int id;
    int a = C17 + 209;
    int b = C52 * 7;
    int c = (209 + C17) % 10000;
    boolean flag = false;
    string label = "rec209";
};
type R210 record {
    int id;
    int a = C18 + 210;
    int b = C55 * 1;
    int c = (210 + C18) % 10000;
    boolean flag = true;
    string label = "rec210";
};
type R211 record {
    int id;
    int a = C19 + 211;
    int b = C58 * 2;
    int c = (211 + C19) % 10000;
    boolean flag = false;
    string label = "rec211";
};
type R212 record {
    int id;
    int a = C20 + 212;
    int b = C61 * 3;
    int c = (212 + C20) % 10000;
    boolean flag = true;
    string label = "rec212";
};
type R213 record {
    int id;
    int a = C21 + 213;
    int b = C0 * 4;
    int c = (213 + C21) % 10000;
    boolean flag = false;
    string label = "rec213";
};
type R214 record {
    int id;
    int a = C22 + 214;
    int b = C3 * 5;
    int c = (214 + C22) % 10000;
    boolean flag = true;
    string label = "rec214";
};
type R215 record {
    int id;
    int a = C23 + 215;
    int b = C6 * 6;
    int c = (215 + C23) % 10000;
    boolean flag = false;
    string label = "rec215";
};
type R216 record {
    int id;
    int a = C24 + 216;
    int b = C9 * 7;
    int c = (216 + C24) % 10000;
    boolean flag = true;
    string label = "rec216";
};
type R217 record {
    int id;
    int a = C25 + 217;
    int b = C12 * 1;
    int c = (217 + C25) % 10000;
    boolean flag = false;
    string label = "rec217";
};
type R218 record {
    int id;
    int a = C26 + 218;
    int b = C15 * 2;
    int c = (218 + C26) % 10000;
    boolean flag = true;
    string label = "rec218";
};
type R219 record {
    int id;
    int a = C27 + 219;
    int b = C18 * 3;
    int c = (219 + C27) % 10000;
    boolean flag = false;
    string label = "rec219";
};
type R220 record {
    int id;
    int a = C28 + 220;
    int b = C21 * 4;
    int c = (220 + C28) % 10000;
    boolean flag = true;
    string label = "rec220";
};
type R221 record {
    int id;
    int a = C29 + 221;
    int b = C24 * 5;
    int c = (221 + C29) % 10000;
    boolean flag = false;
    string label = "rec221";
};
type R222 record {
    int id;
    int a = C30 + 222;
    int b = C27 * 6;
    int c = (222 + C30) % 10000;
    boolean flag = true;
    string label = "rec222";
};
type R223 record {
    int id;
    int a = C31 + 223;
    int b = C30 * 7;
    int c = (223 + C31) % 10000;
    boolean flag = false;
    string label = "rec223";
};
type R224 record {
    int id;
    int a = C32 + 224;
    int b = C33 * 1;
    int c = (224 + C32) % 10000;
    boolean flag = true;
    string label = "rec224";
};
type R225 record {
    int id;
    int a = C33 + 225;
    int b = C36 * 2;
    int c = (225 + C33) % 10000;
    boolean flag = false;
    string label = "rec225";
};
type R226 record {
    int id;
    int a = C34 + 226;
    int b = C39 * 3;
    int c = (226 + C34) % 10000;
    boolean flag = true;
    string label = "rec226";
};
type R227 record {
    int id;
    int a = C35 + 227;
    int b = C42 * 4;
    int c = (227 + C35) % 10000;
    boolean flag = false;
    string label = "rec227";
};
type R228 record {
    int id;
    int a = C36 + 228;
    int b = C45 * 5;
    int c = (228 + C36) % 10000;
    boolean flag = true;
    string label = "rec228";
};
type R229 record {
    int id;
    int a = C37 + 229;
    int b = C48 * 6;
    int c = (229 + C37) % 10000;
    boolean flag = false;
    string label = "rec229";
};
type R230 record {
    int id;
    int a = C38 + 230;
    int b = C51 * 7;
    int c = (230 + C38) % 10000;
    boolean flag = true;
    string label = "rec230";
};
type R231 record {
    int id;
    int a = C39 + 231;
    int b = C54 * 1;
    int c = (231 + C39) % 10000;
    boolean flag = false;
    string label = "rec231";
};
type R232 record {
    int id;
    int a = C40 + 232;
    int b = C57 * 2;
    int c = (232 + C40) % 10000;
    boolean flag = true;
    string label = "rec232";
};
type R233 record {
    int id;
    int a = C41 + 233;
    int b = C60 * 3;
    int c = (233 + C41) % 10000;
    boolean flag = false;
    string label = "rec233";
};
type R234 record {
    int id;
    int a = C42 + 234;
    int b = C63 * 4;
    int c = (234 + C42) % 10000;
    boolean flag = true;
    string label = "rec234";
};
type R235 record {
    int id;
    int a = C43 + 235;
    int b = C2 * 5;
    int c = (235 + C43) % 10000;
    boolean flag = false;
    string label = "rec235";
};
type R236 record {
    int id;
    int a = C44 + 236;
    int b = C5 * 6;
    int c = (236 + C44) % 10000;
    boolean flag = true;
    string label = "rec236";
};
type R237 record {
    int id;
    int a = C45 + 237;
    int b = C8 * 7;
    int c = (237 + C45) % 10000;
    boolean flag = false;
    string label = "rec237";
};
type R238 record {
    int id;
    int a = C46 + 238;
    int b = C11 * 1;
    int c = (238 + C46) % 10000;
    boolean flag = true;
    string label = "rec238";
};
type R239 record {
    int id;
    int a = C47 + 239;
    int b = C14 * 2;
    int c = (239 + C47) % 10000;
    boolean flag = false;
    string label = "rec239";
};
type R240 record {
    int id;
    int a = C48 + 240;
    int b = C17 * 3;
    int c = (240 + C48) % 10000;
    boolean flag = true;
    string label = "rec240";
};
type R241 record {
    int id;
    int a = C49 + 241;
    int b = C20 * 4;
    int c = (241 + C49) % 10000;
    boolean flag = false;
    string label = "rec241";
};
type R242 record {
    int id;
    int a = C50 + 242;
    int b = C23 * 5;
    int c = (242 + C50) % 10000;
    boolean flag = true;
    string label = "rec242";
};
type R243 record {
    int id;
    int a = C51 + 243;
    int b = C26 * 6;
    int c = (243 + C51) % 10000;
    boolean flag = false;
    string label = "rec243";
};
type R244 record {
    int id;
    int a = C52 + 244;
    int b = C29 * 7;
    int c = (244 + C52) % 10000;
    boolean flag = true;
    string label = "rec244";
};
type R245 record {
    int id;
    int a = C53 + 245;
    int b = C32 * 1;
    int c = (245 + C53) % 10000;
    boolean flag = false;
    string label = "rec245";
};
type R246 record {
    int id;
    int a = C54 + 246;
    int b = C35 * 2;
    int c = (246 + C54) % 10000;
    boolean flag = true;
    string label = "rec246";
};
type R247 record {
    int id;
    int a = C55 + 247;
    int b = C38 * 3;
    int c = (247 + C55) % 10000;
    boolean flag = false;
    string label = "rec247";
};
type R248 record {
    int id;
    int a = C56 + 248;
    int b = C41 * 4;
    int c = (248 + C56) % 10000;
    boolean flag = true;
    string label = "rec248";
};
type R249 record {
    int id;
    int a = C57 + 249;
    int b = C44 * 5;
    int c = (249 + C57) % 10000;
    boolean flag = false;
    string label = "rec249";
};
type R250 record {
    int id;
    int a = C58 + 250;
    int b = C47 * 6;
    int c = (250 + C58) % 10000;
    boolean flag = true;
    string label = "rec250";
};
type R251 record {
    int id;
    int a = C59 + 251;
    int b = C50 * 7;
    int c = (251 + C59) % 10000;
    boolean flag = false;
    string label = "rec251";
};
type R252 record {
    int id;
    int a = C60 + 252;
    int b = C53 * 1;
    int c = (252 + C60) % 10000;
    boolean flag = true;
    string label = "rec252";
};
type R253 record {
    int id;
    int a = C61 + 253;
    int b = C56 * 2;
    int c = (253 + C61) % 10000;
    boolean flag = false;
    string label = "rec253";
};
type R254 record {
    int id;
    int a = C62 + 254;
    int b = C59 * 3;
    int c = (254 + C62) % 10000;
    boolean flag = true;
    string label = "rec254";
};
type R255 record {
    int id;
    int a = C63 + 255;
    int b = C62 * 4;
    int c = (255 + C63) % 10000;
    boolean flag = false;
    string label = "rec255";
};
type R256 record {
    int id;
    int a = C0 + 256;
    int b = C1 * 5;
    int c = (256 + C0) % 10000;
    boolean flag = true;
    string label = "rec256";
};
type R257 record {
    int id;
    int a = C1 + 257;
    int b = C4 * 6;
    int c = (257 + C1) % 10000;
    boolean flag = false;
    string label = "rec257";
};
type R258 record {
    int id;
    int a = C2 + 258;
    int b = C7 * 7;
    int c = (258 + C2) % 10000;
    boolean flag = true;
    string label = "rec258";
};
type R259 record {
    int id;
    int a = C3 + 259;
    int b = C10 * 1;
    int c = (259 + C3) % 10000;
    boolean flag = false;
    string label = "rec259";
};
type R260 record {
    int id;
    int a = C4 + 260;
    int b = C13 * 2;
    int c = (260 + C4) % 10000;
    boolean flag = true;
    string label = "rec260";
};
type R261 record {
    int id;
    int a = C5 + 261;
    int b = C16 * 3;
    int c = (261 + C5) % 10000;
    boolean flag = false;
    string label = "rec261";
};
type R262 record {
    int id;
    int a = C6 + 262;
    int b = C19 * 4;
    int c = (262 + C6) % 10000;
    boolean flag = true;
    string label = "rec262";
};
type R263 record {
    int id;
    int a = C7 + 263;
    int b = C22 * 5;
    int c = (263 + C7) % 10000;
    boolean flag = false;
    string label = "rec263";
};
type R264 record {
    int id;
    int a = C8 + 264;
    int b = C25 * 6;
    int c = (264 + C8) % 10000;
    boolean flag = true;
    string label = "rec264";
};
type R265 record {
    int id;
    int a = C9 + 265;
    int b = C28 * 7;
    int c = (265 + C9) % 10000;
    boolean flag = false;
    string label = "rec265";
};
type R266 record {
    int id;
    int a = C10 + 266;
    int b = C31 * 1;
    int c = (266 + C10) % 10000;
    boolean flag = true;
    string label = "rec266";
};
type R267 record {
    int id;
    int a = C11 + 267;
    int b = C34 * 2;
    int c = (267 + C11) % 10000;
    boolean flag = false;
    string label = "rec267";
};
type R268 record {
    int id;
    int a = C12 + 268;
    int b = C37 * 3;
    int c = (268 + C12) % 10000;
    boolean flag = true;
    string label = "rec268";
};
type R269 record {
    int id;
    int a = C13 + 269;
    int b = C40 * 4;
    int c = (269 + C13) % 10000;
    boolean flag = false;
    string label = "rec269";
};
type R270 record {
    int id;
    int a = C14 + 270;
    int b = C43 * 5;
    int c = (270 + C14) % 10000;
    boolean flag = true;
    string label = "rec270";
};
type R271 record {
    int id;
    int a = C15 + 271;
    int b = C46 * 6;
    int c = (271 + C15) % 10000;
    boolean flag = false;
    string label = "rec271";
};
type R272 record {
    int id;
    int a = C16 + 272;
    int b = C49 * 7;
    int c = (272 + C16) % 10000;
    boolean flag = true;
    string label = "rec272";
};
type R273 record {
    int id;
    int a = C17 + 273;
    int b = C52 * 1;
    int c = (273 + C17) % 10000;
    boolean flag = false;
    string label = "rec273";
};
type R274 record {
    int id;
    int a = C18 + 274;
    int b = C55 * 2;
    int c = (274 + C18) % 10000;
    boolean flag = true;
    string label = "rec274";
};
type R275 record {
    int id;
    int a = C19 + 275;
    int b = C58 * 3;
    int c = (275 + C19) % 10000;
    boolean flag = false;
    string label = "rec275";
};
type R276 record {
    int id;
    int a = C20 + 276;
    int b = C61 * 4;
    int c = (276 + C20) % 10000;
    boolean flag = true;
    string label = "rec276";
};
type R277 record {
    int id;
    int a = C21 + 277;
    int b = C0 * 5;
    int c = (277 + C21) % 10000;
    boolean flag = false;
    string label = "rec277";
};
type R278 record {
    int id;
    int a = C22 + 278;
    int b = C3 * 6;
    int c = (278 + C22) % 10000;
    boolean flag = true;
    string label = "rec278";
};
type R279 record {
    int id;
    int a = C23 + 279;
    int b = C6 * 7;
    int c = (279 + C23) % 10000;
    boolean flag = false;
    string label = "rec279";
};
type R280 record {
    int id;
    int a = C24 + 280;
    int b = C9 * 1;
    int c = (280 + C24) % 10000;
    boolean flag = true;
    string label = "rec280";
};
type R281 record {
    int id;
    int a = C25 + 281;
    int b = C12 * 2;
    int c = (281 + C25) % 10000;
    boolean flag = false;
    string label = "rec281";
};
type R282 record {
    int id;
    int a = C26 + 282;
    int b = C15 * 3;
    int c = (282 + C26) % 10000;
    boolean flag = true;
    string label = "rec282";
};
type R283 record {
    int id;
    int a = C27 + 283;
    int b = C18 * 4;
    int c = (283 + C27) % 10000;
    boolean flag = false;
    string label = "rec283";
};
type R284 record {
    int id;
    int a = C28 + 284;
    int b = C21 * 5;
    int c = (284 + C28) % 10000;
    boolean flag = true;
    string label = "rec284";
};
type R285 record {
    int id;
    int a = C29 + 285;
    int b = C24 * 6;
    int c = (285 + C29) % 10000;
    boolean flag = false;
    string label = "rec285";
};
type R286 record {
    int id;
    int a = C30 + 286;
    int b = C27 * 7;
    int c = (286 + C30) % 10000;
    boolean flag = true;
    string label = "rec286";
};
type R287 record {
    int id;
    int a = C31 + 287;
    int b = C30 * 1;
    int c = (287 + C31) % 10000;
    boolean flag = false;
    string label = "rec287";
};
type R288 record {
    int id;
    int a = C32 + 288;
    int b = C33 * 2;
    int c = (288 + C32) % 10000;
    boolean flag = true;
    string label = "rec288";
};
type R289 record {
    int id;
    int a = C33 + 289;
    int b = C36 * 3;
    int c = (289 + C33) % 10000;
    boolean flag = false;
    string label = "rec289";
};
type R290 record {
    int id;
    int a = C34 + 290;
    int b = C39 * 4;
    int c = (290 + C34) % 10000;
    boolean flag = true;
    string label = "rec290";
};
type R291 record {
    int id;
    int a = C35 + 291;
    int b = C42 * 5;
    int c = (291 + C35) % 10000;
    boolean flag = false;
    string label = "rec291";
};
type R292 record {
    int id;
    int a = C36 + 292;
    int b = C45 * 6;
    int c = (292 + C36) % 10000;
    boolean flag = true;
    string label = "rec292";
};
type R293 record {
    int id;
    int a = C37 + 293;
    int b = C48 * 7;
    int c = (293 + C37) % 10000;
    boolean flag = false;
    string label = "rec293";
};
type R294 record {
    int id;
    int a = C38 + 294;
    int b = C51 * 1;
    int c = (294 + C38) % 10000;
    boolean flag = true;
    string label = "rec294";
};
type R295 record {
    int id;
    int a = C39 + 295;
    int b = C54 * 2;
    int c = (295 + C39) % 10000;
    boolean flag = false;
    string label = "rec295";
};
type R296 record {
    int id;
    int a = C40 + 296;
    int b = C57 * 3;
    int c = (296 + C40) % 10000;
    boolean flag = true;
    string label = "rec296";
};
type R297 record {
    int id;
    int a = C41 + 297;
    int b = C60 * 4;
    int c = (297 + C41) % 10000;
    boolean flag = false;
    string label = "rec297";
};
type R298 record {
    int id;
    int a = C42 + 298;
    int b = C63 * 5;
    int c = (298 + C42) % 10000;
    boolean flag = true;
    string label = "rec298";
};
type R299 record {
    int id;
    int a = C43 + 299;
    int b = C2 * 6;
    int c = (299 + C43) % 10000;
    boolean flag = false;
    string label = "rec299";
};
type R300 record {
    int id;
    int a = C44 + 300;
    int b = C5 * 7;
    int c = (300 + C44) % 10000;
    boolean flag = true;
    string label = "rec300";
};
type R301 record {
    int id;
    int a = C45 + 301;
    int b = C8 * 1;
    int c = (301 + C45) % 10000;
    boolean flag = false;
    string label = "rec301";
};
type R302 record {
    int id;
    int a = C46 + 302;
    int b = C11 * 2;
    int c = (302 + C46) % 10000;
    boolean flag = true;
    string label = "rec302";
};
type R303 record {
    int id;
    int a = C47 + 303;
    int b = C14 * 3;
    int c = (303 + C47) % 10000;
    boolean flag = false;
    string label = "rec303";
};
type R304 record {
    int id;
    int a = C48 + 304;
    int b = C17 * 4;
    int c = (304 + C48) % 10000;
    boolean flag = true;
    string label = "rec304";
};
type R305 record {
    int id;
    int a = C49 + 305;
    int b = C20 * 5;
    int c = (305 + C49) % 10000;
    boolean flag = false;
    string label = "rec305";
};
type R306 record {
    int id;
    int a = C50 + 306;
    int b = C23 * 6;
    int c = (306 + C50) % 10000;
    boolean flag = true;
    string label = "rec306";
};
type R307 record {
    int id;
    int a = C51 + 307;
    int b = C26 * 7;
    int c = (307 + C51) % 10000;
    boolean flag = false;
    string label = "rec307";
};
type R308 record {
    int id;
    int a = C52 + 308;
    int b = C29 * 1;
    int c = (308 + C52) % 10000;
    boolean flag = true;
    string label = "rec308";
};
type R309 record {
    int id;
    int a = C53 + 309;
    int b = C32 * 2;
    int c = (309 + C53) % 10000;
    boolean flag = false;
    string label = "rec309";
};
type R310 record {
    int id;
    int a = C54 + 310;
    int b = C35 * 3;
    int c = (310 + C54) % 10000;
    boolean flag = true;
    string label = "rec310";
};
type R311 record {
    int id;
    int a = C55 + 311;
    int b = C38 * 4;
    int c = (311 + C55) % 10000;
    boolean flag = false;
    string label = "rec311";
};
type R312 record {
    int id;
    int a = C56 + 312;
    int b = C41 * 5;
    int c = (312 + C56) % 10000;
    boolean flag = true;
    string label = "rec312";
};
type R313 record {
    int id;
    int a = C57 + 313;
    int b = C44 * 6;
    int c = (313 + C57) % 10000;
    boolean flag = false;
    string label = "rec313";
};
type R314 record {
    int id;
    int a = C58 + 314;
    int b = C47 * 7;
    int c = (314 + C58) % 10000;
    boolean flag = true;
    string label = "rec314";
};
type R315 record {
    int id;
    int a = C59 + 315;
    int b = C50 * 1;
    int c = (315 + C59) % 10000;
    boolean flag = false;
    string label = "rec315";
};
type R316 record {
    int id;
    int a = C60 + 316;
    int b = C53 * 2;
    int c = (316 + C60) % 10000;
    boolean flag = true;
    string label = "rec316";
};
type R317 record {
    int id;
    int a = C61 + 317;
    int b = C56 * 3;
    int c = (317 + C61) % 10000;
    boolean flag = false;
    string label = "rec317";
};
type R318 record {
    int id;
    int a = C62 + 318;
    int b = C59 * 4;
    int c = (318 + C62) % 10000;
    boolean flag = true;
    string label = "rec318";
};
type R319 record {
    int id;
    int a = C63 + 319;
    int b = C62 * 5;
    int c = (319 + C63) % 10000;
    boolean flag = false;
    string label = "rec319";
};
type R320 record {
    int id;
    int a = C0 + 320;
    int b = C1 * 6;
    int c = (320 + C0) % 10000;
    boolean flag = true;
    string label = "rec320";
};
type R321 record {
    int id;
    int a = C1 + 321;
    int b = C4 * 7;
    int c = (321 + C1) % 10000;
    boolean flag = false;
    string label = "rec321";
};
type R322 record {
    int id;
    int a = C2 + 322;
    int b = C7 * 1;
    int c = (322 + C2) % 10000;
    boolean flag = true;
    string label = "rec322";
};
type R323 record {
    int id;
    int a = C3 + 323;
    int b = C10 * 2;
    int c = (323 + C3) % 10000;
    boolean flag = false;
    string label = "rec323";
};
type R324 record {
    int id;
    int a = C4 + 324;
    int b = C13 * 3;
    int c = (324 + C4) % 10000;
    boolean flag = true;
    string label = "rec324";
};
type R325 record {
    int id;
    int a = C5 + 325;
    int b = C16 * 4;
    int c = (325 + C5) % 10000;
    boolean flag = false;
    string label = "rec325";
};
type R326 record {
    int id;
    int a = C6 + 326;
    int b = C19 * 5;
    int c = (326 + C6) % 10000;
    boolean flag = true;
    string label = "rec326";
};
type R327 record {
    int id;
    int a = C7 + 327;
    int b = C22 * 6;
    int c = (327 + C7) % 10000;
    boolean flag = false;
    string label = "rec327";
};
type R328 record {
    int id;
    int a = C8 + 328;
    int b = C25 * 7;
    int c = (328 + C8) % 10000;
    boolean flag = true;
    string label = "rec328";
};
type R329 record {
    int id;
    int a = C9 + 329;
    int b = C28 * 1;
    int c = (329 + C9) % 10000;
    boolean flag = false;
    string label = "rec329";
};
type R330 record {
    int id;
    int a = C10 + 330;
    int b = C31 * 2;
    int c = (330 + C10) % 10000;
    boolean flag = true;
    string label = "rec330";
};
type R331 record {
    int id;
    int a = C11 + 331;
    int b = C34 * 3;
    int c = (331 + C11) % 10000;
    boolean flag = false;
    string label = "rec331";
};
type R332 record {
    int id;
    int a = C12 + 332;
    int b = C37 * 4;
    int c = (332 + C12) % 10000;
    boolean flag = true;
    string label = "rec332";
};
type R333 record {
    int id;
    int a = C13 + 333;
    int b = C40 * 5;
    int c = (333 + C13) % 10000;
    boolean flag = false;
    string label = "rec333";
};
type R334 record {
    int id;
    int a = C14 + 334;
    int b = C43 * 6;
    int c = (334 + C14) % 10000;
    boolean flag = true;
    string label = "rec334";
};
type R335 record {
    int id;
    int a = C15 + 335;
    int b = C46 * 7;
    int c = (335 + C15) % 10000;
    boolean flag = false;
    string label = "rec335";
};
type R336 record {
    int id;
    int a = C16 + 336;
    int b = C49 * 1;
    int c = (336 + C16) % 10000;
    boolean flag = true;
    string label = "rec336";
};
type R337 record {
    int id;
    int a = C17 + 337;
    int b = C52 * 2;
    int c = (337 + C17) % 10000;
    boolean flag = false;
    string label = "rec337";
};
type R338 record {
    int id;
    int a = C18 + 338;
    int b = C55 * 3;
    int c = (338 + C18) % 10000;
    boolean flag = true;
    string label = "rec338";
};
type R339 record {
    int id;
    int a = C19 + 339;
    int b = C58 * 4;
    int c = (339 + C19) % 10000;
    boolean flag = false;
    string label = "rec339";
};
type R340 record {
    int id;
    int a = C20 + 340;
    int b = C61 * 5;
    int c = (340 + C20) % 10000;
    boolean flag = true;
    string label = "rec340";
};
type R341 record {
    int id;
    int a = C21 + 341;
    int b = C0 * 6;
    int c = (341 + C21) % 10000;
    boolean flag = false;
    string label = "rec341";
};
type R342 record {
    int id;
    int a = C22 + 342;
    int b = C3 * 7;
    int c = (342 + C22) % 10000;
    boolean flag = true;
    string label = "rec342";
};
type R343 record {
    int id;
    int a = C23 + 343;
    int b = C6 * 1;
    int c = (343 + C23) % 10000;
    boolean flag = false;
    string label = "rec343";
};
type R344 record {
    int id;
    int a = C24 + 344;
    int b = C9 * 2;
    int c = (344 + C24) % 10000;
    boolean flag = true;
    string label = "rec344";
};
type R345 record {
    int id;
    int a = C25 + 345;
    int b = C12 * 3;
    int c = (345 + C25) % 10000;
    boolean flag = false;
    string label = "rec345";
};
type R346 record {
    int id;
    int a = C26 + 346;
    int b = C15 * 4;
    int c = (346 + C26) % 10000;
    boolean flag = true;
    string label = "rec346";
};
type R347 record {
    int id;
    int a = C27 + 347;
    int b = C18 * 5;
    int c = (347 + C27) % 10000;
    boolean flag = false;
    string label = "rec347";
};
type R348 record {
    int id;
    int a = C28 + 348;
    int b = C21 * 6;
    int c = (348 + C28) % 10000;
    boolean flag = true;
    string label = "rec348";
};
type R349 record {
    int id;
    int a = C29 + 349;
    int b = C24 * 7;
    int c = (349 + C29) % 10000;
    boolean flag = false;
    string label = "rec349";
};
type R350 record {
    int id;
    int a = C30 + 350;
    int b = C27 * 1;
    int c = (350 + C30) % 10000;
    boolean flag = true;
    string label = "rec350";
};
type R351 record {
    int id;
    int a = C31 + 351;
    int b = C30 * 2;
    int c = (351 + C31) % 10000;
    boolean flag = false;
    string label = "rec351";
};
type R352 record {
    int id;
    int a = C32 + 352;
    int b = C33 * 3;
    int c = (352 + C32) % 10000;
    boolean flag = true;
    string label = "rec352";
};
type R353 record {
    int id;
    int a = C33 + 353;
    int b = C36 * 4;
    int c = (353 + C33) % 10000;
    boolean flag = false;
    string label = "rec353";
};
type R354 record {
    int id;
    int a = C34 + 354;
    int b = C39 * 5;
    int c = (354 + C34) % 10000;
    boolean flag = true;
    string label = "rec354";
};
type R355 record {
    int id;
    int a = C35 + 355;
    int b = C42 * 6;
    int c = (355 + C35) % 10000;
    boolean flag = false;
    string label = "rec355";
};
type R356 record {
    int id;
    int a = C36 + 356;
    int b = C45 * 7;
    int c = (356 + C36) % 10000;
    boolean flag = true;
    string label = "rec356";
};
type R357 record {
    int id;
    int a = C37 + 357;
    int b = C48 * 1;
    int c = (357 + C37) % 10000;
    boolean flag = false;
    string label = "rec357";
};
type R358 record {
    int id;
    int a = C38 + 358;
    int b = C51 * 2;
    int c = (358 + C38) % 10000;
    boolean flag = true;
    string label = "rec358";
};
type R359 record {
    int id;
    int a = C39 + 359;
    int b = C54 * 3;
    int c = (359 + C39) % 10000;
    boolean flag = false;
    string label = "rec359";
};
type R360 record {
    int id;
    int a = C40 + 360;
    int b = C57 * 4;
    int c = (360 + C40) % 10000;
    boolean flag = true;
    string label = "rec360";
};
type R361 record {
    int id;
    int a = C41 + 361;
    int b = C60 * 5;
    int c = (361 + C41) % 10000;
    boolean flag = false;
    string label = "rec361";
};
type R362 record {
    int id;
    int a = C42 + 362;
    int b = C63 * 6;
    int c = (362 + C42) % 10000;
    boolean flag = true;
    string label = "rec362";
};
type R363 record {
    int id;
    int a = C43 + 363;
    int b = C2 * 7;
    int c = (363 + C43) % 10000;
    boolean flag = false;
    string label = "rec363";
};
type R364 record {
    int id;
    int a = C44 + 364;
    int b = C5 * 1;
    int c = (364 + C44) % 10000;
    boolean flag = true;
    string label = "rec364";
};
type R365 record {
    int id;
    int a = C45 + 365;
    int b = C8 * 2;
    int c = (365 + C45) % 10000;
    boolean flag = false;
    string label = "rec365";
};
type R366 record {
    int id;
    int a = C46 + 366;
    int b = C11 * 3;
    int c = (366 + C46) % 10000;
    boolean flag = true;
    string label = "rec366";
};
type R367 record {
    int id;
    int a = C47 + 367;
    int b = C14 * 4;
    int c = (367 + C47) % 10000;
    boolean flag = false;
    string label = "rec367";
};
type R368 record {
    int id;
    int a = C48 + 368;
    int b = C17 * 5;
    int c = (368 + C48) % 10000;
    boolean flag = true;
    string label = "rec368";
};
type R369 record {
    int id;
    int a = C49 + 369;
    int b = C20 * 6;
    int c = (369 + C49) % 10000;
    boolean flag = false;
    string label = "rec369";
};
type R370 record {
    int id;
    int a = C50 + 370;
    int b = C23 * 7;
    int c = (370 + C50) % 10000;
    boolean flag = true;
    string label = "rec370";
};
type R371 record {
    int id;
    int a = C51 + 371;
    int b = C26 * 1;
    int c = (371 + C51) % 10000;
    boolean flag = false;
    string label = "rec371";
};
type R372 record {
    int id;
    int a = C52 + 372;
    int b = C29 * 2;
    int c = (372 + C52) % 10000;
    boolean flag = true;
    string label = "rec372";
};
type R373 record {
    int id;
    int a = C53 + 373;
    int b = C32 * 3;
    int c = (373 + C53) % 10000;
    boolean flag = false;
    string label = "rec373";
};
type R374 record {
    int id;
    int a = C54 + 374;
    int b = C35 * 4;
    int c = (374 + C54) % 10000;
    boolean flag = true;
    string label = "rec374";
};
type R375 record {
    int id;
    int a = C55 + 375;
    int b = C38 * 5;
    int c = (375 + C55) % 10000;
    boolean flag = false;
    string label = "rec375";
};
type R376 record {
    int id;
    int a = C56 + 376;
    int b = C41 * 6;
    int c = (376 + C56) % 10000;
    boolean flag = true;
    string label = "rec376";
};
type R377 record {
    int id;
    int a = C57 + 377;
    int b = C44 * 7;
    int c = (377 + C57) % 10000;
    boolean flag = false;
    string label = "rec377";
};
type R378 record {
    int id;
    int a = C58 + 378;
    int b = C47 * 1;
    int c = (378 + C58) % 10000;
    boolean flag = true;
    string label = "rec378";
};
type R379 record {
    int id;
    int a = C59 + 379;
    int b = C50 * 2;
    int c = (379 + C59) % 10000;
    boolean flag = false;
    string label = "rec379";
};
type R380 record {
    int id;
    int a = C60 + 380;
    int b = C53 * 3;
    int c = (380 + C60) % 10000;
    boolean flag = true;
    string label = "rec380";
};
type R381 record {
    int id;
    int a = C61 + 381;
    int b = C56 * 4;
    int c = (381 + C61) % 10000;
    boolean flag = false;
    string label = "rec381";
};
type R382 record {
    int id;
    int a = C62 + 382;
    int b = C59 * 5;
    int c = (382 + C62) % 10000;
    boolean flag = true;
    string label = "rec382";
};
type R383 record {
    int id;
    int a = C63 + 383;
    int b = C62 * 6;
    int c = (383 + C63) % 10000;
    boolean flag = false;
    string label = "rec383";
};
type R384 record {
    int id;
    int a = C0 + 384;
    int b = C1 * 7;
    int c = (384 + C0) % 10000;
    boolean flag = true;
    string label = "rec384";
};
type R385 record {
    int id;
    int a = C1 + 385;
    int b = C4 * 1;
    int c = (385 + C1) % 10000;
    boolean flag = false;
    string label = "rec385";
};
type R386 record {
    int id;
    int a = C2 + 386;
    int b = C7 * 2;
    int c = (386 + C2) % 10000;
    boolean flag = true;
    string label = "rec386";
};
type R387 record {
    int id;
    int a = C3 + 387;
    int b = C10 * 3;
    int c = (387 + C3) % 10000;
    boolean flag = false;
    string label = "rec387";
};
type R388 record {
    int id;
    int a = C4 + 388;
    int b = C13 * 4;
    int c = (388 + C4) % 10000;
    boolean flag = true;
    string label = "rec388";
};
type R389 record {
    int id;
    int a = C5 + 389;
    int b = C16 * 5;
    int c = (389 + C5) % 10000;
    boolean flag = false;
    string label = "rec389";
};
type R390 record {
    int id;
    int a = C6 + 390;
    int b = C19 * 6;
    int c = (390 + C6) % 10000;
    boolean flag = true;
    string label = "rec390";
};
type R391 record {
    int id;
    int a = C7 + 391;
    int b = C22 * 7;
    int c = (391 + C7) % 10000;
    boolean flag = false;
    string label = "rec391";
};
type R392 record {
    int id;
    int a = C8 + 392;
    int b = C25 * 1;
    int c = (392 + C8) % 10000;
    boolean flag = true;
    string label = "rec392";
};
type R393 record {
    int id;
    int a = C9 + 393;
    int b = C28 * 2;
    int c = (393 + C9) % 10000;
    boolean flag = false;
    string label = "rec393";
};
type R394 record {
    int id;
    int a = C10 + 394;
    int b = C31 * 3;
    int c = (394 + C10) % 10000;
    boolean flag = true;
    string label = "rec394";
};
type R395 record {
    int id;
    int a = C11 + 395;
    int b = C34 * 4;
    int c = (395 + C11) % 10000;
    boolean flag = false;
    string label = "rec395";
};
type R396 record {
    int id;
    int a = C12 + 396;
    int b = C37 * 5;
    int c = (396 + C12) % 10000;
    boolean flag = true;
    string label = "rec396";
};
type R397 record {
    int id;
    int a = C13 + 397;
    int b = C40 * 6;
    int c = (397 + C13) % 10000;
    boolean flag = false;
    string label = "rec397";
};
type R398 record {
    int id;
    int a = C14 + 398;
    int b = C43 * 7;
    int c = (398 + C14) % 10000;
    boolean flag = true;
    string label = "rec398";
};
type R399 record {
    int id;
    int a = C15 + 399;
    int b = C46 * 1;
    int c = (399 + C15) % 10000;
    boolean flag = false;
    string label = "rec399";
};
type R400 record {
    int id;
    int a = C16 + 400;
    int b = C49 * 2;
    int c = (400 + C16) % 10000;
    boolean flag = true;
    string label = "rec400";
};
type R401 record {
    int id;
    int a = C17 + 401;
    int b = C52 * 3;
    int c = (401 + C17) % 10000;
    boolean flag = false;
    string label = "rec401";
};
type R402 record {
    int id;
    int a = C18 + 402;
    int b = C55 * 4;
    int c = (402 + C18) % 10000;
    boolean flag = true;
    string label = "rec402";
};
type R403 record {
    int id;
    int a = C19 + 403;
    int b = C58 * 5;
    int c = (403 + C19) % 10000;
    boolean flag = false;
    string label = "rec403";
};
type R404 record {
    int id;
    int a = C20 + 404;
    int b = C61 * 6;
    int c = (404 + C20) % 10000;
    boolean flag = true;
    string label = "rec404";
};
type R405 record {
    int id;
    int a = C21 + 405;
    int b = C0 * 7;
    int c = (405 + C21) % 10000;
    boolean flag = false;
    string label = "rec405";
};
type R406 record {
    int id;
    int a = C22 + 406;
    int b = C3 * 1;
    int c = (406 + C22) % 10000;
    boolean flag = true;
    string label = "rec406";
};
type R407 record {
    int id;
    int a = C23 + 407;
    int b = C6 * 2;
    int c = (407 + C23) % 10000;
    boolean flag = false;
    string label = "rec407";
};
type R408 record {
    int id;
    int a = C24 + 408;
    int b = C9 * 3;
    int c = (408 + C24) % 10000;
    boolean flag = true;
    string label = "rec408";
};
type R409 record {
    int id;
    int a = C25 + 409;
    int b = C12 * 4;
    int c = (409 + C25) % 10000;
    boolean flag = false;
    string label = "rec409";
};
type R410 record {
    int id;
    int a = C26 + 410;
    int b = C15 * 5;
    int c = (410 + C26) % 10000;
    boolean flag = true;
    string label = "rec410";
};
type R411 record {
    int id;
    int a = C27 + 411;
    int b = C18 * 6;
    int c = (411 + C27) % 10000;
    boolean flag = false;
    string label = "rec411";
};
type R412 record {
    int id;
    int a = C28 + 412;
    int b = C21 * 7;
    int c = (412 + C28) % 10000;
    boolean flag = true;
    string label = "rec412";
};
type R413 record {
    int id;
    int a = C29 + 413;
    int b = C24 * 1;
    int c = (413 + C29) % 10000;
    boolean flag = false;
    string label = "rec413";
};
type R414 record {
    int id;
    int a = C30 + 414;
    int b = C27 * 2;
    int c = (414 + C30) % 10000;
    boolean flag = true;
    string label = "rec414";
};
type R415 record {
    int id;
    int a = C31 + 415;
    int b = C30 * 3;
    int c = (415 + C31) % 10000;
    boolean flag = false;
    string label = "rec415";
};
type R416 record {
    int id;
    int a = C32 + 416;
    int b = C33 * 4;
    int c = (416 + C32) % 10000;
    boolean flag = true;
    string label = "rec416";
};
type R417 record {
    int id;
    int a = C33 + 417;
    int b = C36 * 5;
    int c = (417 + C33) % 10000;
    boolean flag = false;
    string label = "rec417";
};
type R418 record {
    int id;
    int a = C34 + 418;
    int b = C39 * 6;
    int c = (418 + C34) % 10000;
    boolean flag = true;
    string label = "rec418";
};
type R419 record {
    int id;
    int a = C35 + 419;
    int b = C42 * 7;
    int c = (419 + C35) % 10000;
    boolean flag = false;
    string label = "rec419";
};
type R420 record {
    int id;
    int a = C36 + 420;
    int b = C45 * 1;
    int c = (420 + C36) % 10000;
    boolean flag = true;
    string label = "rec420";
};
type R421 record {
    int id;
    int a = C37 + 421;
    int b = C48 * 2;
    int c = (421 + C37) % 10000;
    boolean flag = false;
    string label = "rec421";
};
type R422 record {
    int id;
    int a = C38 + 422;
    int b = C51 * 3;
    int c = (422 + C38) % 10000;
    boolean flag = true;
    string label = "rec422";
};
type R423 record {
    int id;
    int a = C39 + 423;
    int b = C54 * 4;
    int c = (423 + C39) % 10000;
    boolean flag = false;
    string label = "rec423";
};
type R424 record {
    int id;
    int a = C40 + 424;
    int b = C57 * 5;
    int c = (424 + C40) % 10000;
    boolean flag = true;
    string label = "rec424";
};
type R425 record {
    int id;
    int a = C41 + 425;
    int b = C60 * 6;
    int c = (425 + C41) % 10000;
    boolean flag = false;
    string label = "rec425";
};
type R426 record {
    int id;
    int a = C42 + 426;
    int b = C63 * 7;
    int c = (426 + C42) % 10000;
    boolean flag = true;
    string label = "rec426";
};
type R427 record {
    int id;
    int a = C43 + 427;
    int b = C2 * 1;
    int c = (427 + C43) % 10000;
    boolean flag = false;
    string label = "rec427";
};
type R428 record {
    int id;
    int a = C44 + 428;
    int b = C5 * 2;
    int c = (428 + C44) % 10000;
    boolean flag = true;
    string label = "rec428";
};
type R429 record {
    int id;
    int a = C45 + 429;
    int b = C8 * 3;
    int c = (429 + C45) % 10000;
    boolean flag = false;
    string label = "rec429";
};
type R430 record {
    int id;
    int a = C46 + 430;
    int b = C11 * 4;
    int c = (430 + C46) % 10000;
    boolean flag = true;
    string label = "rec430";
};
type R431 record {
    int id;
    int a = C47 + 431;
    int b = C14 * 5;
    int c = (431 + C47) % 10000;
    boolean flag = false;
    string label = "rec431";
};
type R432 record {
    int id;
    int a = C48 + 432;
    int b = C17 * 6;
    int c = (432 + C48) % 10000;
    boolean flag = true;
    string label = "rec432";
};
type R433 record {
    int id;
    int a = C49 + 433;
    int b = C20 * 7;
    int c = (433 + C49) % 10000;
    boolean flag = false;
    string label = "rec433";
};
type R434 record {
    int id;
    int a = C50 + 434;
    int b = C23 * 1;
    int c = (434 + C50) % 10000;
    boolean flag = true;
    string label = "rec434";
};
type R435 record {
    int id;
    int a = C51 + 435;
    int b = C26 * 2;
    int c = (435 + C51) % 10000;
    boolean flag = false;
    string label = "rec435";
};
type R436 record {
    int id;
    int a = C52 + 436;
    int b = C29 * 3;
    int c = (436 + C52) % 10000;
    boolean flag = true;
    string label = "rec436";
};
type R437 record {
    int id;
    int a = C53 + 437;
    int b = C32 * 4;
    int c = (437 + C53) % 10000;
    boolean flag = false;
    string label = "rec437";
};
type R438 record {
    int id;
    int a = C54 + 438;
    int b = C35 * 5;
    int c = (438 + C54) % 10000;
    boolean flag = true;
    string label = "rec438";
};
type R439 record {
    int id;
    int a = C55 + 439;
    int b = C38 * 6;
    int c = (439 + C55) % 10000;
    boolean flag = false;
    string label = "rec439";
};
type R440 record {
    int id;
    int a = C56 + 440;
    int b = C41 * 7;
    int c = (440 + C56) % 10000;
    boolean flag = true;
    string label = "rec440";
};
type R441 record {
    int id;
    int a = C57 + 441;
    int b = C44 * 1;
    int c = (441 + C57) % 10000;
    boolean flag = false;
    string label = "rec441";
};
type R442 record {
    int id;
    int a = C58 + 442;
    int b = C47 * 2;
    int c = (442 + C58) % 10000;
    boolean flag = true;
    string label = "rec442";
};
type R443 record {
    int id;
    int a = C59 + 443;
    int b = C50 * 3;
    int c = (443 + C59) % 10000;
    boolean flag = false;
    string label = "rec443";
};
type R444 record {
    int id;
    int a = C60 + 444;
    int b = C53 * 4;
    int c = (444 + C60) % 10000;
    boolean flag = true;
    string label = "rec444";
};
type R445 record {
    int id;
    int a = C61 + 445;
    int b = C56 * 5;
    int c = (445 + C61) % 10000;
    boolean flag = false;
    string label = "rec445";
};
type R446 record {
    int id;
    int a = C62 + 446;
    int b = C59 * 6;
    int c = (446 + C62) % 10000;
    boolean flag = true;
    string label = "rec446";
};
type R447 record {
    int id;
    int a = C63 + 447;
    int b = C62 * 7;
    int c = (447 + C63) % 10000;
    boolean flag = false;
    string label = "rec447";
};
type R448 record {
    int id;
    int a = C0 + 448;
    int b = C1 * 1;
    int c = (448 + C0) % 10000;
    boolean flag = true;
    string label = "rec448";
};
type R449 record {
    int id;
    int a = C1 + 449;
    int b = C4 * 2;
    int c = (449 + C1) % 10000;
    boolean flag = false;
    string label = "rec449";
};
type R450 record {
    int id;
    int a = C2 + 450;
    int b = C7 * 3;
    int c = (450 + C2) % 10000;
    boolean flag = true;
    string label = "rec450";
};
type R451 record {
    int id;
    int a = C3 + 451;
    int b = C10 * 4;
    int c = (451 + C3) % 10000;
    boolean flag = false;
    string label = "rec451";
};
type R452 record {
    int id;
    int a = C4 + 452;
    int b = C13 * 5;
    int c = (452 + C4) % 10000;
    boolean flag = true;
    string label = "rec452";
};
type R453 record {
    int id;
    int a = C5 + 453;
    int b = C16 * 6;
    int c = (453 + C5) % 10000;
    boolean flag = false;
    string label = "rec453";
};
type R454 record {
    int id;
    int a = C6 + 454;
    int b = C19 * 7;
    int c = (454 + C6) % 10000;
    boolean flag = true;
    string label = "rec454";
};
type R455 record {
    int id;
    int a = C7 + 455;
    int b = C22 * 1;
    int c = (455 + C7) % 10000;
    boolean flag = false;
    string label = "rec455";
};
type R456 record {
    int id;
    int a = C8 + 456;
    int b = C25 * 2;
    int c = (456 + C8) % 10000;
    boolean flag = true;
    string label = "rec456";
};
type R457 record {
    int id;
    int a = C9 + 457;
    int b = C28 * 3;
    int c = (457 + C9) % 10000;
    boolean flag = false;
    string label = "rec457";
};
type R458 record {
    int id;
    int a = C10 + 458;
    int b = C31 * 4;
    int c = (458 + C10) % 10000;
    boolean flag = true;
    string label = "rec458";
};
type R459 record {
    int id;
    int a = C11 + 459;
    int b = C34 * 5;
    int c = (459 + C11) % 10000;
    boolean flag = false;
    string label = "rec459";
};
type R460 record {
    int id;
    int a = C12 + 460;
    int b = C37 * 6;
    int c = (460 + C12) % 10000;
    boolean flag = true;
    string label = "rec460";
};
type R461 record {
    int id;
    int a = C13 + 461;
    int b = C40 * 7;
    int c = (461 + C13) % 10000;
    boolean flag = false;
    string label = "rec461";
};
type R462 record {
    int id;
    int a = C14 + 462;
    int b = C43 * 1;
    int c = (462 + C14) % 10000;
    boolean flag = true;
    string label = "rec462";
};
type R463 record {
    int id;
    int a = C15 + 463;
    int b = C46 * 2;
    int c = (463 + C15) % 10000;
    boolean flag = false;
    string label = "rec463";
};
type R464 record {
    int id;
    int a = C16 + 464;
    int b = C49 * 3;
    int c = (464 + C16) % 10000;
    boolean flag = true;
    string label = "rec464";
};
type R465 record {
    int id;
    int a = C17 + 465;
    int b = C52 * 4;
    int c = (465 + C17) % 10000;
    boolean flag = false;
    string label = "rec465";
};
type R466 record {
    int id;
    int a = C18 + 466;
    int b = C55 * 5;
    int c = (466 + C18) % 10000;
    boolean flag = true;
    string label = "rec466";
};
type R467 record {
    int id;
    int a = C19 + 467;
    int b = C58 * 6;
    int c = (467 + C19) % 10000;
    boolean flag = false;
    string label = "rec467";
};
type R468 record {
    int id;
    int a = C20 + 468;
    int b = C61 * 7;
    int c = (468 + C20) % 10000;
    boolean flag = true;
    string label = "rec468";
};
type R469 record {
    int id;
    int a = C21 + 469;
    int b = C0 * 1;
    int c = (469 + C21) % 10000;
    boolean flag = false;
    string label = "rec469";
};
type R470 record {
    int id;
    int a = C22 + 470;
    int b = C3 * 2;
    int c = (470 + C22) % 10000;
    boolean flag = true;
    string label = "rec470";
};
type R471 record {
    int id;
    int a = C23 + 471;
    int b = C6 * 3;
    int c = (471 + C23) % 10000;
    boolean flag = false;
    string label = "rec471";
};
type R472 record {
    int id;
    int a = C24 + 472;
    int b = C9 * 4;
    int c = (472 + C24) % 10000;
    boolean flag = true;
    string label = "rec472";
};
type R473 record {
    int id;
    int a = C25 + 473;
    int b = C12 * 5;
    int c = (473 + C25) % 10000;
    boolean flag = false;
    string label = "rec473";
};
type R474 record {
    int id;
    int a = C26 + 474;
    int b = C15 * 6;
    int c = (474 + C26) % 10000;
    boolean flag = true;
    string label = "rec474";
};
type R475 record {
    int id;
    int a = C27 + 475;
    int b = C18 * 7;
    int c = (475 + C27) % 10000;
    boolean flag = false;
    string label = "rec475";
};
type R476 record {
    int id;
    int a = C28 + 476;
    int b = C21 * 1;
    int c = (476 + C28) % 10000;
    boolean flag = true;
    string label = "rec476";
};
type R477 record {
    int id;
    int a = C29 + 477;
    int b = C24 * 2;
    int c = (477 + C29) % 10000;
    boolean flag = false;
    string label = "rec477";
};
type R478 record {
    int id;
    int a = C30 + 478;
    int b = C27 * 3;
    int c = (478 + C30) % 10000;
    boolean flag = true;
    string label = "rec478";
};
type R479 record {
    int id;
    int a = C31 + 479;
    int b = C30 * 4;
    int c = (479 + C31) % 10000;
    boolean flag = false;
    string label = "rec479";
};
type R480 record {
    int id;
    int a = C32 + 480;
    int b = C33 * 5;
    int c = (480 + C32) % 10000;
    boolean flag = true;
    string label = "rec480";
};
type R481 record {
    int id;
    int a = C33 + 481;
    int b = C36 * 6;
    int c = (481 + C33) % 10000;
    boolean flag = false;
    string label = "rec481";
};
type R482 record {
    int id;
    int a = C34 + 482;
    int b = C39 * 7;
    int c = (482 + C34) % 10000;
    boolean flag = true;
    string label = "rec482";
};
type R483 record {
    int id;
    int a = C35 + 483;
    int b = C42 * 1;
    int c = (483 + C35) % 10000;
    boolean flag = false;
    string label = "rec483";
};
type R484 record {
    int id;
    int a = C36 + 484;
    int b = C45 * 2;
    int c = (484 + C36) % 10000;
    boolean flag = true;
    string label = "rec484";
};
type R485 record {
    int id;
    int a = C37 + 485;
    int b = C48 * 3;
    int c = (485 + C37) % 10000;
    boolean flag = false;
    string label = "rec485";
};
type R486 record {
    int id;
    int a = C38 + 486;
    int b = C51 * 4;
    int c = (486 + C38) % 10000;
    boolean flag = true;
    string label = "rec486";
};
type R487 record {
    int id;
    int a = C39 + 487;
    int b = C54 * 5;
    int c = (487 + C39) % 10000;
    boolean flag = false;
    string label = "rec487";
};
type R488 record {
    int id;
    int a = C40 + 488;
    int b = C57 * 6;
    int c = (488 + C40) % 10000;
    boolean flag = true;
    string label = "rec488";
};
type R489 record {
    int id;
    int a = C41 + 489;
    int b = C60 * 7;
    int c = (489 + C41) % 10000;
    boolean flag = false;
    string label = "rec489";
};
type R490 record {
    int id;
    int a = C42 + 490;
    int b = C63 * 1;
    int c = (490 + C42) % 10000;
    boolean flag = true;
    string label = "rec490";
};
type R491 record {
    int id;
    int a = C43 + 491;
    int b = C2 * 2;
    int c = (491 + C43) % 10000;
    boolean flag = false;
    string label = "rec491";
};
type R492 record {
    int id;
    int a = C44 + 492;
    int b = C5 * 3;
    int c = (492 + C44) % 10000;
    boolean flag = true;
    string label = "rec492";
};
type R493 record {
    int id;
    int a = C45 + 493;
    int b = C8 * 4;
    int c = (493 + C45) % 10000;
    boolean flag = false;
    string label = "rec493";
};
type R494 record {
    int id;
    int a = C46 + 494;
    int b = C11 * 5;
    int c = (494 + C46) % 10000;
    boolean flag = true;
    string label = "rec494";
};
type R495 record {
    int id;
    int a = C47 + 495;
    int b = C14 * 6;
    int c = (495 + C47) % 10000;
    boolean flag = false;
    string label = "rec495";
};
type R496 record {
    int id;
    int a = C48 + 496;
    int b = C17 * 7;
    int c = (496 + C48) % 10000;
    boolean flag = true;
    string label = "rec496";
};
type R497 record {
    int id;
    int a = C49 + 497;
    int b = C20 * 1;
    int c = (497 + C49) % 10000;
    boolean flag = false;
    string label = "rec497";
};
type R498 record {
    int id;
    int a = C50 + 498;
    int b = C23 * 2;
    int c = (498 + C50) % 10000;
    boolean flag = true;
    string label = "rec498";
};
type R499 record {
    int id;
    int a = C51 + 499;
    int b = C26 * 3;
    int c = (499 + C51) % 10000;
    boolean flag = false;
    string label = "rec499";
};
type R500 record {
    int id;
    int a = C52 + 500;
    int b = C29 * 4;
    int c = (500 + C52) % 10000;
    boolean flag = true;
    string label = "rec500";
};
type R501 record {
    int id;
    int a = C53 + 501;
    int b = C32 * 5;
    int c = (501 + C53) % 10000;
    boolean flag = false;
    string label = "rec501";
};
type R502 record {
    int id;
    int a = C54 + 502;
    int b = C35 * 6;
    int c = (502 + C54) % 10000;
    boolean flag = true;
    string label = "rec502";
};
type R503 record {
    int id;
    int a = C55 + 503;
    int b = C38 * 7;
    int c = (503 + C55) % 10000;
    boolean flag = false;
    string label = "rec503";
};
type R504 record {
    int id;
    int a = C56 + 504;
    int b = C41 * 1;
    int c = (504 + C56) % 10000;
    boolean flag = true;
    string label = "rec504";
};
type R505 record {
    int id;
    int a = C57 + 505;
    int b = C44 * 2;
    int c = (505 + C57) % 10000;
    boolean flag = false;
    string label = "rec505";
};
type R506 record {
    int id;
    int a = C58 + 506;
    int b = C47 * 3;
    int c = (506 + C58) % 10000;
    boolean flag = true;
    string label = "rec506";
};
type R507 record {
    int id;
    int a = C59 + 507;
    int b = C50 * 4;
    int c = (507 + C59) % 10000;
    boolean flag = false;
    string label = "rec507";
};
type R508 record {
    int id;
    int a = C60 + 508;
    int b = C53 * 5;
    int c = (508 + C60) % 10000;
    boolean flag = true;
    string label = "rec508";
};
type R509 record {
    int id;
    int a = C61 + 509;
    int b = C56 * 6;
    int c = (509 + C61) % 10000;
    boolean flag = false;
    string label = "rec509";
};
type R510 record {
    int id;
    int a = C62 + 510;
    int b = C59 * 7;
    int c = (510 + C62) % 10000;
    boolean flag = true;
    string label = "rec510";
};
type R511 record {
    int id;
    int a = C63 + 511;
    int b = C62 * 1;
    int c = (511 + C63) % 10000;
    boolean flag = false;
    string label = "rec511";
};
type R512 record {
    int id;
    int a = C0 + 512;
    int b = C1 * 2;
    int c = (512 + C0) % 10000;
    boolean flag = true;
    string label = "rec512";
};
type R513 record {
    int id;
    int a = C1 + 513;
    int b = C4 * 3;
    int c = (513 + C1) % 10000;
    boolean flag = false;
    string label = "rec513";
};
type R514 record {
    int id;
    int a = C2 + 514;
    int b = C7 * 4;
    int c = (514 + C2) % 10000;
    boolean flag = true;
    string label = "rec514";
};
type R515 record {
    int id;
    int a = C3 + 515;
    int b = C10 * 5;
    int c = (515 + C3) % 10000;
    boolean flag = false;
    string label = "rec515";
};
type R516 record {
    int id;
    int a = C4 + 516;
    int b = C13 * 6;
    int c = (516 + C4) % 10000;
    boolean flag = true;
    string label = "rec516";
};
type R517 record {
    int id;
    int a = C5 + 517;
    int b = C16 * 7;
    int c = (517 + C5) % 10000;
    boolean flag = false;
    string label = "rec517";
};
type R518 record {
    int id;
    int a = C6 + 518;
    int b = C19 * 1;
    int c = (518 + C6) % 10000;
    boolean flag = true;
    string label = "rec518";
};
type R519 record {
    int id;
    int a = C7 + 519;
    int b = C22 * 2;
    int c = (519 + C7) % 10000;
    boolean flag = false;
    string label = "rec519";
};
type R520 record {
    int id;
    int a = C8 + 520;
    int b = C25 * 3;
    int c = (520 + C8) % 10000;
    boolean flag = true;
    string label = "rec520";
};
type R521 record {
    int id;
    int a = C9 + 521;
    int b = C28 * 4;
    int c = (521 + C9) % 10000;
    boolean flag = false;
    string label = "rec521";
};
type R522 record {
    int id;
    int a = C10 + 522;
    int b = C31 * 5;
    int c = (522 + C10) % 10000;
    boolean flag = true;
    string label = "rec522";
};
type R523 record {
    int id;
    int a = C11 + 523;
    int b = C34 * 6;
    int c = (523 + C11) % 10000;
    boolean flag = false;
    string label = "rec523";
};
type R524 record {
    int id;
    int a = C12 + 524;
    int b = C37 * 7;
    int c = (524 + C12) % 10000;
    boolean flag = true;
    string label = "rec524";
};
type R525 record {
    int id;
    int a = C13 + 525;
    int b = C40 * 1;
    int c = (525 + C13) % 10000;
    boolean flag = false;
    string label = "rec525";
};
type R526 record {
    int id;
    int a = C14 + 526;
    int b = C43 * 2;
    int c = (526 + C14) % 10000;
    boolean flag = true;
    string label = "rec526";
};
type R527 record {
    int id;
    int a = C15 + 527;
    int b = C46 * 3;
    int c = (527 + C15) % 10000;
    boolean flag = false;
    string label = "rec527";
};
type R528 record {
    int id;
    int a = C16 + 528;
    int b = C49 * 4;
    int c = (528 + C16) % 10000;
    boolean flag = true;
    string label = "rec528";
};
type R529 record {
    int id;
    int a = C17 + 529;
    int b = C52 * 5;
    int c = (529 + C17) % 10000;
    boolean flag = false;
    string label = "rec529";
};
type R530 record {
    int id;
    int a = C18 + 530;
    int b = C55 * 6;
    int c = (530 + C18) % 10000;
    boolean flag = true;
    string label = "rec530";
};
type R531 record {
    int id;
    int a = C19 + 531;
    int b = C58 * 7;
    int c = (531 + C19) % 10000;
    boolean flag = false;
    string label = "rec531";
};
type R532 record {
    int id;
    int a = C20 + 532;
    int b = C61 * 1;
    int c = (532 + C20) % 10000;
    boolean flag = true;
    string label = "rec532";
};
type R533 record {
    int id;
    int a = C21 + 533;
    int b = C0 * 2;
    int c = (533 + C21) % 10000;
    boolean flag = false;
    string label = "rec533";
};
type R534 record {
    int id;
    int a = C22 + 534;
    int b = C3 * 3;
    int c = (534 + C22) % 10000;
    boolean flag = true;
    string label = "rec534";
};
type R535 record {
    int id;
    int a = C23 + 535;
    int b = C6 * 4;
    int c = (535 + C23) % 10000;
    boolean flag = false;
    string label = "rec535";
};
type R536 record {
    int id;
    int a = C24 + 536;
    int b = C9 * 5;
    int c = (536 + C24) % 10000;
    boolean flag = true;
    string label = "rec536";
};
type R537 record {
    int id;
    int a = C25 + 537;
    int b = C12 * 6;
    int c = (537 + C25) % 10000;
    boolean flag = false;
    string label = "rec537";
};
type R538 record {
    int id;
    int a = C26 + 538;
    int b = C15 * 7;
    int c = (538 + C26) % 10000;
    boolean flag = true;
    string label = "rec538";
};
type R539 record {
    int id;
    int a = C27 + 539;
    int b = C18 * 1;
    int c = (539 + C27) % 10000;
    boolean flag = false;
    string label = "rec539";
};
type R540 record {
    int id;
    int a = C28 + 540;
    int b = C21 * 2;
    int c = (540 + C28) % 10000;
    boolean flag = true;
    string label = "rec540";
};
type R541 record {
    int id;
    int a = C29 + 541;
    int b = C24 * 3;
    int c = (541 + C29) % 10000;
    boolean flag = false;
    string label = "rec541";
};
type R542 record {
    int id;
    int a = C30 + 542;
    int b = C27 * 4;
    int c = (542 + C30) % 10000;
    boolean flag = true;
    string label = "rec542";
};
type R543 record {
    int id;
    int a = C31 + 543;
    int b = C30 * 5;
    int c = (543 + C31) % 10000;
    boolean flag = false;
    string label = "rec543";
};
type R544 record {
    int id;
    int a = C32 + 544;
    int b = C33 * 6;
    int c = (544 + C32) % 10000;
    boolean flag = true;
    string label = "rec544";
};
type R545 record {
    int id;
    int a = C33 + 545;
    int b = C36 * 7;
    int c = (545 + C33) % 10000;
    boolean flag = false;
    string label = "rec545";
};
type R546 record {
    int id;
    int a = C34 + 546;
    int b = C39 * 1;
    int c = (546 + C34) % 10000;
    boolean flag = true;
    string label = "rec546";
};
type R547 record {
    int id;
    int a = C35 + 547;
    int b = C42 * 2;
    int c = (547 + C35) % 10000;
    boolean flag = false;
    string label = "rec547";
};
type R548 record {
    int id;
    int a = C36 + 548;
    int b = C45 * 3;
    int c = (548 + C36) % 10000;
    boolean flag = true;
    string label = "rec548";
};
type R549 record {
    int id;
    int a = C37 + 549;
    int b = C48 * 4;
    int c = (549 + C37) % 10000;
    boolean flag = false;
    string label = "rec549";
};
type R550 record {
    int id;
    int a = C38 + 550;
    int b = C51 * 5;
    int c = (550 + C38) % 10000;
    boolean flag = true;
    string label = "rec550";
};
type R551 record {
    int id;
    int a = C39 + 551;
    int b = C54 * 6;
    int c = (551 + C39) % 10000;
    boolean flag = false;
    string label = "rec551";
};
type R552 record {
    int id;
    int a = C40 + 552;
    int b = C57 * 7;
    int c = (552 + C40) % 10000;
    boolean flag = true;
    string label = "rec552";
};
type R553 record {
    int id;
    int a = C41 + 553;
    int b = C60 * 1;
    int c = (553 + C41) % 10000;
    boolean flag = false;
    string label = "rec553";
};
type R554 record {
    int id;
    int a = C42 + 554;
    int b = C63 * 2;
    int c = (554 + C42) % 10000;
    boolean flag = true;
    string label = "rec554";
};
type R555 record {
    int id;
    int a = C43 + 555;
    int b = C2 * 3;
    int c = (555 + C43) % 10000;
    boolean flag = false;
    string label = "rec555";
};
type R556 record {
    int id;
    int a = C44 + 556;
    int b = C5 * 4;
    int c = (556 + C44) % 10000;
    boolean flag = true;
    string label = "rec556";
};
type R557 record {
    int id;
    int a = C45 + 557;
    int b = C8 * 5;
    int c = (557 + C45) % 10000;
    boolean flag = false;
    string label = "rec557";
};
type R558 record {
    int id;
    int a = C46 + 558;
    int b = C11 * 6;
    int c = (558 + C46) % 10000;
    boolean flag = true;
    string label = "rec558";
};
type R559 record {
    int id;
    int a = C47 + 559;
    int b = C14 * 7;
    int c = (559 + C47) % 10000;
    boolean flag = false;
    string label = "rec559";
};
type R560 record {
    int id;
    int a = C48 + 560;
    int b = C17 * 1;
    int c = (560 + C48) % 10000;
    boolean flag = true;
    string label = "rec560";
};
type R561 record {
    int id;
    int a = C49 + 561;
    int b = C20 * 2;
    int c = (561 + C49) % 10000;
    boolean flag = false;
    string label = "rec561";
};
type R562 record {
    int id;
    int a = C50 + 562;
    int b = C23 * 3;
    int c = (562 + C50) % 10000;
    boolean flag = true;
    string label = "rec562";
};
type R563 record {
    int id;
    int a = C51 + 563;
    int b = C26 * 4;
    int c = (563 + C51) % 10000;
    boolean flag = false;
    string label = "rec563";
};
type R564 record {
    int id;
    int a = C52 + 564;
    int b = C29 * 5;
    int c = (564 + C52) % 10000;
    boolean flag = true;
    string label = "rec564";
};
type R565 record {
    int id;
    int a = C53 + 565;
    int b = C32 * 6;
    int c = (565 + C53) % 10000;
    boolean flag = false;
    string label = "rec565";
};
type R566 record {
    int id;
    int a = C54 + 566;
    int b = C35 * 7;
    int c = (566 + C54) % 10000;
    boolean flag = true;
    string label = "rec566";
};
type R567 record {
    int id;
    int a = C55 + 567;
    int b = C38 * 1;
    int c = (567 + C55) % 10000;
    boolean flag = false;
    string label = "rec567";
};
type R568 record {
    int id;
    int a = C56 + 568;
    int b = C41 * 2;
    int c = (568 + C56) % 10000;
    boolean flag = true;
    string label = "rec568";
};
type R569 record {
    int id;
    int a = C57 + 569;
    int b = C44 * 3;
    int c = (569 + C57) % 10000;
    boolean flag = false;
    string label = "rec569";
};
type R570 record {
    int id;
    int a = C58 + 570;
    int b = C47 * 4;
    int c = (570 + C58) % 10000;
    boolean flag = true;
    string label = "rec570";
};
type R571 record {
    int id;
    int a = C59 + 571;
    int b = C50 * 5;
    int c = (571 + C59) % 10000;
    boolean flag = false;
    string label = "rec571";
};
type R572 record {
    int id;
    int a = C60 + 572;
    int b = C53 * 6;
    int c = (572 + C60) % 10000;
    boolean flag = true;
    string label = "rec572";
};
type R573 record {
    int id;
    int a = C61 + 573;
    int b = C56 * 7;
    int c = (573 + C61) % 10000;
    boolean flag = false;
    string label = "rec573";
};
type R574 record {
    int id;
    int a = C62 + 574;
    int b = C59 * 1;
    int c = (574 + C62) % 10000;
    boolean flag = true;
    string label = "rec574";
};
type R575 record {
    int id;
    int a = C63 + 575;
    int b = C62 * 2;
    int c = (575 + C63) % 10000;
    boolean flag = false;
    string label = "rec575";
};
type R576 record {
    int id;
    int a = C0 + 576;
    int b = C1 * 3;
    int c = (576 + C0) % 10000;
    boolean flag = true;
    string label = "rec576";
};
type R577 record {
    int id;
    int a = C1 + 577;
    int b = C4 * 4;
    int c = (577 + C1) % 10000;
    boolean flag = false;
    string label = "rec577";
};
type R578 record {
    int id;
    int a = C2 + 578;
    int b = C7 * 5;
    int c = (578 + C2) % 10000;
    boolean flag = true;
    string label = "rec578";
};
type R579 record {
    int id;
    int a = C3 + 579;
    int b = C10 * 6;
    int c = (579 + C3) % 10000;
    boolean flag = false;
    string label = "rec579";
};
type R580 record {
    int id;
    int a = C4 + 580;
    int b = C13 * 7;
    int c = (580 + C4) % 10000;
    boolean flag = true;
    string label = "rec580";
};
type R581 record {
    int id;
    int a = C5 + 581;
    int b = C16 * 1;
    int c = (581 + C5) % 10000;
    boolean flag = false;
    string label = "rec581";
};
type R582 record {
    int id;
    int a = C6 + 582;
    int b = C19 * 2;
    int c = (582 + C6) % 10000;
    boolean flag = true;
    string label = "rec582";
};
type R583 record {
    int id;
    int a = C7 + 583;
    int b = C22 * 3;
    int c = (583 + C7) % 10000;
    boolean flag = false;
    string label = "rec583";
};
type R584 record {
    int id;
    int a = C8 + 584;
    int b = C25 * 4;
    int c = (584 + C8) % 10000;
    boolean flag = true;
    string label = "rec584";
};
type R585 record {
    int id;
    int a = C9 + 585;
    int b = C28 * 5;
    int c = (585 + C9) % 10000;
    boolean flag = false;
    string label = "rec585";
};
type R586 record {
    int id;
    int a = C10 + 586;
    int b = C31 * 6;
    int c = (586 + C10) % 10000;
    boolean flag = true;
    string label = "rec586";
};
type R587 record {
    int id;
    int a = C11 + 587;
    int b = C34 * 7;
    int c = (587 + C11) % 10000;
    boolean flag = false;
    string label = "rec587";
};
type R588 record {
    int id;
    int a = C12 + 588;
    int b = C37 * 1;
    int c = (588 + C12) % 10000;
    boolean flag = true;
    string label = "rec588";
};
type R589 record {
    int id;
    int a = C13 + 589;
    int b = C40 * 2;
    int c = (589 + C13) % 10000;
    boolean flag = false;
    string label = "rec589";
};
type R590 record {
    int id;
    int a = C14 + 590;
    int b = C43 * 3;
    int c = (590 + C14) % 10000;
    boolean flag = true;
    string label = "rec590";
};
type R591 record {
    int id;
    int a = C15 + 591;
    int b = C46 * 4;
    int c = (591 + C15) % 10000;
    boolean flag = false;
    string label = "rec591";
};
type R592 record {
    int id;
    int a = C16 + 592;
    int b = C49 * 5;
    int c = (592 + C16) % 10000;
    boolean flag = true;
    string label = "rec592";
};
type R593 record {
    int id;
    int a = C17 + 593;
    int b = C52 * 6;
    int c = (593 + C17) % 10000;
    boolean flag = false;
    string label = "rec593";
};
type R594 record {
    int id;
    int a = C18 + 594;
    int b = C55 * 7;
    int c = (594 + C18) % 10000;
    boolean flag = true;
    string label = "rec594";
};
type R595 record {
    int id;
    int a = C19 + 595;
    int b = C58 * 1;
    int c = (595 + C19) % 10000;
    boolean flag = false;
    string label = "rec595";
};
type R596 record {
    int id;
    int a = C20 + 596;
    int b = C61 * 2;
    int c = (596 + C20) % 10000;
    boolean flag = true;
    string label = "rec596";
};
type R597 record {
    int id;
    int a = C21 + 597;
    int b = C0 * 3;
    int c = (597 + C21) % 10000;
    boolean flag = false;
    string label = "rec597";
};
type R598 record {
    int id;
    int a = C22 + 598;
    int b = C3 * 4;
    int c = (598 + C22) % 10000;
    boolean flag = true;
    string label = "rec598";
};
type R599 record {
    int id;
    int a = C23 + 599;
    int b = C6 * 5;
    int c = (599 + C23) % 10000;
    boolean flag = false;
    string label = "rec599";
};
type R600 record {
    int id;
    int a = C24 + 600;
    int b = C9 * 6;
    int c = (600 + C24) % 10000;
    boolean flag = true;
    string label = "rec600";
};
type R601 record {
    int id;
    int a = C25 + 601;
    int b = C12 * 7;
    int c = (601 + C25) % 10000;
    boolean flag = false;
    string label = "rec601";
};
type R602 record {
    int id;
    int a = C26 + 602;
    int b = C15 * 1;
    int c = (602 + C26) % 10000;
    boolean flag = true;
    string label = "rec602";
};
type R603 record {
    int id;
    int a = C27 + 603;
    int b = C18 * 2;
    int c = (603 + C27) % 10000;
    boolean flag = false;
    string label = "rec603";
};
type R604 record {
    int id;
    int a = C28 + 604;
    int b = C21 * 3;
    int c = (604 + C28) % 10000;
    boolean flag = true;
    string label = "rec604";
};
type R605 record {
    int id;
    int a = C29 + 605;
    int b = C24 * 4;
    int c = (605 + C29) % 10000;
    boolean flag = false;
    string label = "rec605";
};
type R606 record {
    int id;
    int a = C30 + 606;
    int b = C27 * 5;
    int c = (606 + C30) % 10000;
    boolean flag = true;
    string label = "rec606";
};
type R607 record {
    int id;
    int a = C31 + 607;
    int b = C30 * 6;
    int c = (607 + C31) % 10000;
    boolean flag = false;
    string label = "rec607";
};
type R608 record {
    int id;
    int a = C32 + 608;
    int b = C33 * 7;
    int c = (608 + C32) % 10000;
    boolean flag = true;
    string label = "rec608";
};
type R609 record {
    int id;
    int a = C33 + 609;
    int b = C36 * 1;
    int c = (609 + C33) % 10000;
    boolean flag = false;
    string label = "rec609";
};
type R610 record {
    int id;
    int a = C34 + 610;
    int b = C39 * 2;
    int c = (610 + C34) % 10000;
    boolean flag = true;
    string label = "rec610";
};
type R611 record {
    int id;
    int a = C35 + 611;
    int b = C42 * 3;
    int c = (611 + C35) % 10000;
    boolean flag = false;
    string label = "rec611";
};
type R612 record {
    int id;
    int a = C36 + 612;
    int b = C45 * 4;
    int c = (612 + C36) % 10000;
    boolean flag = true;
    string label = "rec612";
};
type R613 record {
    int id;
    int a = C37 + 613;
    int b = C48 * 5;
    int c = (613 + C37) % 10000;
    boolean flag = false;
    string label = "rec613";
};
type R614 record {
    int id;
    int a = C38 + 614;
    int b = C51 * 6;
    int c = (614 + C38) % 10000;
    boolean flag = true;
    string label = "rec614";
};
type R615 record {
    int id;
    int a = C39 + 615;
    int b = C54 * 7;
    int c = (615 + C39) % 10000;
    boolean flag = false;
    string label = "rec615";
};
type R616 record {
    int id;
    int a = C40 + 616;
    int b = C57 * 1;
    int c = (616 + C40) % 10000;
    boolean flag = true;
    string label = "rec616";
};
type R617 record {
    int id;
    int a = C41 + 617;
    int b = C60 * 2;
    int c = (617 + C41) % 10000;
    boolean flag = false;
    string label = "rec617";
};
type R618 record {
    int id;
    int a = C42 + 618;
    int b = C63 * 3;
    int c = (618 + C42) % 10000;
    boolean flag = true;
    string label = "rec618";
};
type R619 record {
    int id;
    int a = C43 + 619;
    int b = C2 * 4;
    int c = (619 + C43) % 10000;
    boolean flag = false;
    string label = "rec619";
};
type R620 record {
    int id;
    int a = C44 + 620;
    int b = C5 * 5;
    int c = (620 + C44) % 10000;
    boolean flag = true;
    string label = "rec620";
};
type R621 record {
    int id;
    int a = C45 + 621;
    int b = C8 * 6;
    int c = (621 + C45) % 10000;
    boolean flag = false;
    string label = "rec621";
};
type R622 record {
    int id;
    int a = C46 + 622;
    int b = C11 * 7;
    int c = (622 + C46) % 10000;
    boolean flag = true;
    string label = "rec622";
};
type R623 record {
    int id;
    int a = C47 + 623;
    int b = C14 * 1;
    int c = (623 + C47) % 10000;
    boolean flag = false;
    string label = "rec623";
};
type R624 record {
    int id;
    int a = C48 + 624;
    int b = C17 * 2;
    int c = (624 + C48) % 10000;
    boolean flag = true;
    string label = "rec624";
};
type R625 record {
    int id;
    int a = C49 + 625;
    int b = C20 * 3;
    int c = (625 + C49) % 10000;
    boolean flag = false;
    string label = "rec625";
};
type R626 record {
    int id;
    int a = C50 + 626;
    int b = C23 * 4;
    int c = (626 + C50) % 10000;
    boolean flag = true;
    string label = "rec626";
};
type R627 record {
    int id;
    int a = C51 + 627;
    int b = C26 * 5;
    int c = (627 + C51) % 10000;
    boolean flag = false;
    string label = "rec627";
};
type R628 record {
    int id;
    int a = C52 + 628;
    int b = C29 * 6;
    int c = (628 + C52) % 10000;
    boolean flag = true;
    string label = "rec628";
};
type R629 record {
    int id;
    int a = C53 + 629;
    int b = C32 * 7;
    int c = (629 + C53) % 10000;
    boolean flag = false;
    string label = "rec629";
};
type R630 record {
    int id;
    int a = C54 + 630;
    int b = C35 * 1;
    int c = (630 + C54) % 10000;
    boolean flag = true;
    string label = "rec630";
};
type R631 record {
    int id;
    int a = C55 + 631;
    int b = C38 * 2;
    int c = (631 + C55) % 10000;
    boolean flag = false;
    string label = "rec631";
};
type R632 record {
    int id;
    int a = C56 + 632;
    int b = C41 * 3;
    int c = (632 + C56) % 10000;
    boolean flag = true;
    string label = "rec632";
};
type R633 record {
    int id;
    int a = C57 + 633;
    int b = C44 * 4;
    int c = (633 + C57) % 10000;
    boolean flag = false;
    string label = "rec633";
};
type R634 record {
    int id;
    int a = C58 + 634;
    int b = C47 * 5;
    int c = (634 + C58) % 10000;
    boolean flag = true;
    string label = "rec634";
};
type R635 record {
    int id;
    int a = C59 + 635;
    int b = C50 * 6;
    int c = (635 + C59) % 10000;
    boolean flag = false;
    string label = "rec635";
};
type R636 record {
    int id;
    int a = C60 + 636;
    int b = C53 * 7;
    int c = (636 + C60) % 10000;
    boolean flag = true;
    string label = "rec636";
};
type R637 record {
    int id;
    int a = C61 + 637;
    int b = C56 * 1;
    int c = (637 + C61) % 10000;
    boolean flag = false;
    string label = "rec637";
};
type R638 record {
    int id;
    int a = C62 + 638;
    int b = C59 * 2;
    int c = (638 + C62) % 10000;
    boolean flag = true;
    string label = "rec638";
};
type R639 record {
    int id;
    int a = C63 + 639;
    int b = C62 * 3;
    int c = (639 + C63) % 10000;
    boolean flag = false;
    string label = "rec639";
};
type R640 record {
    int id;
    int a = C0 + 640;
    int b = C1 * 4;
    int c = (640 + C0) % 10000;
    boolean flag = true;
    string label = "rec640";
};
type R641 record {
    int id;
    int a = C1 + 641;
    int b = C4 * 5;
    int c = (641 + C1) % 10000;
    boolean flag = false;
    string label = "rec641";
};
type R642 record {
    int id;
    int a = C2 + 642;
    int b = C7 * 6;
    int c = (642 + C2) % 10000;
    boolean flag = true;
    string label = "rec642";
};
type R643 record {
    int id;
    int a = C3 + 643;
    int b = C10 * 7;
    int c = (643 + C3) % 10000;
    boolean flag = false;
    string label = "rec643";
};
type R644 record {
    int id;
    int a = C4 + 644;
    int b = C13 * 1;
    int c = (644 + C4) % 10000;
    boolean flag = true;
    string label = "rec644";
};
type R645 record {
    int id;
    int a = C5 + 645;
    int b = C16 * 2;
    int c = (645 + C5) % 10000;
    boolean flag = false;
    string label = "rec645";
};
type R646 record {
    int id;
    int a = C6 + 646;
    int b = C19 * 3;
    int c = (646 + C6) % 10000;
    boolean flag = true;
    string label = "rec646";
};
type R647 record {
    int id;
    int a = C7 + 647;
    int b = C22 * 4;
    int c = (647 + C7) % 10000;
    boolean flag = false;
    string label = "rec647";
};
type R648 record {
    int id;
    int a = C8 + 648;
    int b = C25 * 5;
    int c = (648 + C8) % 10000;
    boolean flag = true;
    string label = "rec648";
};
type R649 record {
    int id;
    int a = C9 + 649;
    int b = C28 * 6;
    int c = (649 + C9) % 10000;
    boolean flag = false;
    string label = "rec649";
};
type R650 record {
    int id;
    int a = C10 + 650;
    int b = C31 * 7;
    int c = (650 + C10) % 10000;
    boolean flag = true;
    string label = "rec650";
};
type R651 record {
    int id;
    int a = C11 + 651;
    int b = C34 * 1;
    int c = (651 + C11) % 10000;
    boolean flag = false;
    string label = "rec651";
};
type R652 record {
    int id;
    int a = C12 + 652;
    int b = C37 * 2;
    int c = (652 + C12) % 10000;
    boolean flag = true;
    string label = "rec652";
};
type R653 record {
    int id;
    int a = C13 + 653;
    int b = C40 * 3;
    int c = (653 + C13) % 10000;
    boolean flag = false;
    string label = "rec653";
};
type R654 record {
    int id;
    int a = C14 + 654;
    int b = C43 * 4;
    int c = (654 + C14) % 10000;
    boolean flag = true;
    string label = "rec654";
};
type R655 record {
    int id;
    int a = C15 + 655;
    int b = C46 * 5;
    int c = (655 + C15) % 10000;
    boolean flag = false;
    string label = "rec655";
};
type R656 record {
    int id;
    int a = C16 + 656;
    int b = C49 * 6;
    int c = (656 + C16) % 10000;
    boolean flag = true;
    string label = "rec656";
};
type R657 record {
    int id;
    int a = C17 + 657;
    int b = C52 * 7;
    int c = (657 + C17) % 10000;
    boolean flag = false;
    string label = "rec657";
};
type R658 record {
    int id;
    int a = C18 + 658;
    int b = C55 * 1;
    int c = (658 + C18) % 10000;
    boolean flag = true;
    string label = "rec658";
};
type R659 record {
    int id;
    int a = C19 + 659;
    int b = C58 * 2;
    int c = (659 + C19) % 10000;
    boolean flag = false;
    string label = "rec659";
};
type R660 record {
    int id;
    int a = C20 + 660;
    int b = C61 * 3;
    int c = (660 + C20) % 10000;
    boolean flag = true;
    string label = "rec660";
};
type R661 record {
    int id;
    int a = C21 + 661;
    int b = C0 * 4;
    int c = (661 + C21) % 10000;
    boolean flag = false;
    string label = "rec661";
};
type R662 record {
    int id;
    int a = C22 + 662;
    int b = C3 * 5;
    int c = (662 + C22) % 10000;
    boolean flag = true;
    string label = "rec662";
};
type R663 record {
    int id;
    int a = C23 + 663;
    int b = C6 * 6;
    int c = (663 + C23) % 10000;
    boolean flag = false;
    string label = "rec663";
};
type R664 record {
    int id;
    int a = C24 + 664;
    int b = C9 * 7;
    int c = (664 + C24) % 10000;
    boolean flag = true;
    string label = "rec664";
};
type R665 record {
    int id;
    int a = C25 + 665;
    int b = C12 * 1;
    int c = (665 + C25) % 10000;
    boolean flag = false;
    string label = "rec665";
};
type R666 record {
    int id;
    int a = C26 + 666;
    int b = C15 * 2;
    int c = (666 + C26) % 10000;
    boolean flag = true;
    string label = "rec666";
};
type R667 record {
    int id;
    int a = C27 + 667;
    int b = C18 * 3;
    int c = (667 + C27) % 10000;
    boolean flag = false;
    string label = "rec667";
};
type R668 record {
    int id;
    int a = C28 + 668;
    int b = C21 * 4;
    int c = (668 + C28) % 10000;
    boolean flag = true;
    string label = "rec668";
};
type R669 record {
    int id;
    int a = C29 + 669;
    int b = C24 * 5;
    int c = (669 + C29) % 10000;
    boolean flag = false;
    string label = "rec669";
};
type R670 record {
    int id;
    int a = C30 + 670;
    int b = C27 * 6;
    int c = (670 + C30) % 10000;
    boolean flag = true;
    string label = "rec670";
};
type R671 record {
    int id;
    int a = C31 + 671;
    int b = C30 * 7;
    int c = (671 + C31) % 10000;
    boolean flag = false;
    string label = "rec671";
};
type R672 record {
    int id;
    int a = C32 + 672;
    int b = C33 * 1;
    int c = (672 + C32) % 10000;
    boolean flag = true;
    string label = "rec672";
};
type R673 record {
    int id;
    int a = C33 + 673;
    int b = C36 * 2;
    int c = (673 + C33) % 10000;
    boolean flag = false;
    string label = "rec673";
};
type R674 record {
    int id;
    int a = C34 + 674;
    int b = C39 * 3;
    int c = (674 + C34) % 10000;
    boolean flag = true;
    string label = "rec674";
};
type R675 record {
    int id;
    int a = C35 + 675;
    int b = C42 * 4;
    int c = (675 + C35) % 10000;
    boolean flag = false;
    string label = "rec675";
};
type R676 record {
    int id;
    int a = C36 + 676;
    int b = C45 * 5;
    int c = (676 + C36) % 10000;
    boolean flag = true;
    string label = "rec676";
};
type R677 record {
    int id;
    int a = C37 + 677;
    int b = C48 * 6;
    int c = (677 + C37) % 10000;
    boolean flag = false;
    string label = "rec677";
};
type R678 record {
    int id;
    int a = C38 + 678;
    int b = C51 * 7;
    int c = (678 + C38) % 10000;
    boolean flag = true;
    string label = "rec678";
};
type R679 record {
    int id;
    int a = C39 + 679;
    int b = C54 * 1;
    int c = (679 + C39) % 10000;
    boolean flag = false;
    string label = "rec679";
};
type R680 record {
    int id;
    int a = C40 + 680;
    int b = C57 * 2;
    int c = (680 + C40) % 10000;
    boolean flag = true;
    string label = "rec680";
};
type R681 record {
    int id;
    int a = C41 + 681;
    int b = C60 * 3;
    int c = (681 + C41) % 10000;
    boolean flag = false;
    string label = "rec681";
};
type R682 record {
    int id;
    int a = C42 + 682;
    int b = C63 * 4;
    int c = (682 + C42) % 10000;
    boolean flag = true;
    string label = "rec682";
};
type R683 record {
    int id;
    int a = C43 + 683;
    int b = C2 * 5;
    int c = (683 + C43) % 10000;
    boolean flag = false;
    string label = "rec683";
};
type R684 record {
    int id;
    int a = C44 + 684;
    int b = C5 * 6;
    int c = (684 + C44) % 10000;
    boolean flag = true;
    string label = "rec684";
};
type R685 record {
    int id;
    int a = C45 + 685;
    int b = C8 * 7;
    int c = (685 + C45) % 10000;
    boolean flag = false;
    string label = "rec685";
};
type R686 record {
    int id;
    int a = C46 + 686;
    int b = C11 * 1;
    int c = (686 + C46) % 10000;
    boolean flag = true;
    string label = "rec686";
};
type R687 record {
    int id;
    int a = C47 + 687;
    int b = C14 * 2;
    int c = (687 + C47) % 10000;
    boolean flag = false;
    string label = "rec687";
};
type R688 record {
    int id;
    int a = C48 + 688;
    int b = C17 * 3;
    int c = (688 + C48) % 10000;
    boolean flag = true;
    string label = "rec688";
};
type R689 record {
    int id;
    int a = C49 + 689;
    int b = C20 * 4;
    int c = (689 + C49) % 10000;
    boolean flag = false;
    string label = "rec689";
};
type R690 record {
    int id;
    int a = C50 + 690;
    int b = C23 * 5;
    int c = (690 + C50) % 10000;
    boolean flag = true;
    string label = "rec690";
};
type R691 record {
    int id;
    int a = C51 + 691;
    int b = C26 * 6;
    int c = (691 + C51) % 10000;
    boolean flag = false;
    string label = "rec691";
};
type R692 record {
    int id;
    int a = C52 + 692;
    int b = C29 * 7;
    int c = (692 + C52) % 10000;
    boolean flag = true;
    string label = "rec692";
};
type R693 record {
    int id;
    int a = C53 + 693;
    int b = C32 * 1;
    int c = (693 + C53) % 10000;
    boolean flag = false;
    string label = "rec693";
};
type R694 record {
    int id;
    int a = C54 + 694;
    int b = C35 * 2;
    int c = (694 + C54) % 10000;
    boolean flag = true;
    string label = "rec694";
};
type R695 record {
    int id;
    int a = C55 + 695;
    int b = C38 * 3;
    int c = (695 + C55) % 10000;
    boolean flag = false;
    string label = "rec695";
};
type R696 record {
    int id;
    int a = C56 + 696;
    int b = C41 * 4;
    int c = (696 + C56) % 10000;
    boolean flag = true;
    string label = "rec696";
};
type R697 record {
    int id;
    int a = C57 + 697;
    int b = C44 * 5;
    int c = (697 + C57) % 10000;
    boolean flag = false;
    string label = "rec697";
};
type R698 record {
    int id;
    int a = C58 + 698;
    int b = C47 * 6;
    int c = (698 + C58) % 10000;
    boolean flag = true;
    string label = "rec698";
};
type R699 record {
    int id;
    int a = C59 + 699;
    int b = C50 * 7;
    int c = (699 + C59) % 10000;
    boolean flag = false;
    string label = "rec699";
};
type R700 record {
    int id;
    int a = C60 + 700;
    int b = C53 * 1;
    int c = (700 + C60) % 10000;
    boolean flag = true;
    string label = "rec700";
};
type R701 record {
    int id;
    int a = C61 + 701;
    int b = C56 * 2;
    int c = (701 + C61) % 10000;
    boolean flag = false;
    string label = "rec701";
};
type R702 record {
    int id;
    int a = C62 + 702;
    int b = C59 * 3;
    int c = (702 + C62) % 10000;
    boolean flag = true;
    string label = "rec702";
};
type R703 record {
    int id;
    int a = C63 + 703;
    int b = C62 * 4;
    int c = (703 + C63) % 10000;
    boolean flag = false;
    string label = "rec703";
};
type R704 record {
    int id;
    int a = C0 + 704;
    int b = C1 * 5;
    int c = (704 + C0) % 10000;
    boolean flag = true;
    string label = "rec704";
};
type R705 record {
    int id;
    int a = C1 + 705;
    int b = C4 * 6;
    int c = (705 + C1) % 10000;
    boolean flag = false;
    string label = "rec705";
};
type R706 record {
    int id;
    int a = C2 + 706;
    int b = C7 * 7;
    int c = (706 + C2) % 10000;
    boolean flag = true;
    string label = "rec706";
};
type R707 record {
    int id;
    int a = C3 + 707;
    int b = C10 * 1;
    int c = (707 + C3) % 10000;
    boolean flag = false;
    string label = "rec707";
};
type R708 record {
    int id;
    int a = C4 + 708;
    int b = C13 * 2;
    int c = (708 + C4) % 10000;
    boolean flag = true;
    string label = "rec708";
};
type R709 record {
    int id;
    int a = C5 + 709;
    int b = C16 * 3;
    int c = (709 + C5) % 10000;
    boolean flag = false;
    string label = "rec709";
};
type R710 record {
    int id;
    int a = C6 + 710;
    int b = C19 * 4;
    int c = (710 + C6) % 10000;
    boolean flag = true;
    string label = "rec710";
};
type R711 record {
    int id;
    int a = C7 + 711;
    int b = C22 * 5;
    int c = (711 + C7) % 10000;
    boolean flag = false;
    string label = "rec711";
};
type R712 record {
    int id;
    int a = C8 + 712;
    int b = C25 * 6;
    int c = (712 + C8) % 10000;
    boolean flag = true;
    string label = "rec712";
};
type R713 record {
    int id;
    int a = C9 + 713;
    int b = C28 * 7;
    int c = (713 + C9) % 10000;
    boolean flag = false;
    string label = "rec713";
};
type R714 record {
    int id;
    int a = C10 + 714;
    int b = C31 * 1;
    int c = (714 + C10) % 10000;
    boolean flag = true;
    string label = "rec714";
};
type R715 record {
    int id;
    int a = C11 + 715;
    int b = C34 * 2;
    int c = (715 + C11) % 10000;
    boolean flag = false;
    string label = "rec715";
};
type R716 record {
    int id;
    int a = C12 + 716;
    int b = C37 * 3;
    int c = (716 + C12) % 10000;
    boolean flag = true;
    string label = "rec716";
};
type R717 record {
    int id;
    int a = C13 + 717;
    int b = C40 * 4;
    int c = (717 + C13) % 10000;
    boolean flag = false;
    string label = "rec717";
};
type R718 record {
    int id;
    int a = C14 + 718;
    int b = C43 * 5;
    int c = (718 + C14) % 10000;
    boolean flag = true;
    string label = "rec718";
};
type R719 record {
    int id;
    int a = C15 + 719;
    int b = C46 * 6;
    int c = (719 + C15) % 10000;
    boolean flag = false;
    string label = "rec719";
};
type R720 record {
    int id;
    int a = C16 + 720;
    int b = C49 * 7;
    int c = (720 + C16) % 10000;
    boolean flag = true;
    string label = "rec720";
};
type R721 record {
    int id;
    int a = C17 + 721;
    int b = C52 * 1;
    int c = (721 + C17) % 10000;
    boolean flag = false;
    string label = "rec721";
};
type R722 record {
    int id;
    int a = C18 + 722;
    int b = C55 * 2;
    int c = (722 + C18) % 10000;
    boolean flag = true;
    string label = "rec722";
};
type R723 record {
    int id;
    int a = C19 + 723;
    int b = C58 * 3;
    int c = (723 + C19) % 10000;
    boolean flag = false;
    string label = "rec723";
};
type R724 record {
    int id;
    int a = C20 + 724;
    int b = C61 * 4;
    int c = (724 + C20) % 10000;
    boolean flag = true;
    string label = "rec724";
};
type R725 record {
    int id;
    int a = C21 + 725;
    int b = C0 * 5;
    int c = (725 + C21) % 10000;
    boolean flag = false;
    string label = "rec725";
};
type R726 record {
    int id;
    int a = C22 + 726;
    int b = C3 * 6;
    int c = (726 + C22) % 10000;
    boolean flag = true;
    string label = "rec726";
};
type R727 record {
    int id;
    int a = C23 + 727;
    int b = C6 * 7;
    int c = (727 + C23) % 10000;
    boolean flag = false;
    string label = "rec727";
};
type R728 record {
    int id;
    int a = C24 + 728;
    int b = C9 * 1;
    int c = (728 + C24) % 10000;
    boolean flag = true;
    string label = "rec728";
};
type R729 record {
    int id;
    int a = C25 + 729;
    int b = C12 * 2;
    int c = (729 + C25) % 10000;
    boolean flag = false;
    string label = "rec729";
};
type R730 record {
    int id;
    int a = C26 + 730;
    int b = C15 * 3;
    int c = (730 + C26) % 10000;
    boolean flag = true;
    string label = "rec730";
};
type R731 record {
    int id;
    int a = C27 + 731;
    int b = C18 * 4;
    int c = (731 + C27) % 10000;
    boolean flag = false;
    string label = "rec731";
};
type R732 record {
    int id;
    int a = C28 + 732;
    int b = C21 * 5;
    int c = (732 + C28) % 10000;
    boolean flag = true;
    string label = "rec732";
};
type R733 record {
    int id;
    int a = C29 + 733;
    int b = C24 * 6;
    int c = (733 + C29) % 10000;
    boolean flag = false;
    string label = "rec733";
};
type R734 record {
    int id;
    int a = C30 + 734;
    int b = C27 * 7;
    int c = (734 + C30) % 10000;
    boolean flag = true;
    string label = "rec734";
};
type R735 record {
    int id;
    int a = C31 + 735;
    int b = C30 * 1;
    int c = (735 + C31) % 10000;
    boolean flag = false;
    string label = "rec735";
};
type R736 record {
    int id;
    int a = C32 + 736;
    int b = C33 * 2;
    int c = (736 + C32) % 10000;
    boolean flag = true;
    string label = "rec736";
};
type R737 record {
    int id;
    int a = C33 + 737;
    int b = C36 * 3;
    int c = (737 + C33) % 10000;
    boolean flag = false;
    string label = "rec737";
};
type R738 record {
    int id;
    int a = C34 + 738;
    int b = C39 * 4;
    int c = (738 + C34) % 10000;
    boolean flag = true;
    string label = "rec738";
};
type R739 record {
    int id;
    int a = C35 + 739;
    int b = C42 * 5;
    int c = (739 + C35) % 10000;
    boolean flag = false;
    string label = "rec739";
};
type R740 record {
    int id;
    int a = C36 + 740;
    int b = C45 * 6;
    int c = (740 + C36) % 10000;
    boolean flag = true;
    string label = "rec740";
};
type R741 record {
    int id;
    int a = C37 + 741;
    int b = C48 * 7;
    int c = (741 + C37) % 10000;
    boolean flag = false;
    string label = "rec741";
};
type R742 record {
    int id;
    int a = C38 + 742;
    int b = C51 * 1;
    int c = (742 + C38) % 10000;
    boolean flag = true;
    string label = "rec742";
};
type R743 record {
    int id;
    int a = C39 + 743;
    int b = C54 * 2;
    int c = (743 + C39) % 10000;
    boolean flag = false;
    string label = "rec743";
};
type R744 record {
    int id;
    int a = C40 + 744;
    int b = C57 * 3;
    int c = (744 + C40) % 10000;
    boolean flag = true;
    string label = "rec744";
};
type R745 record {
    int id;
    int a = C41 + 745;
    int b = C60 * 4;
    int c = (745 + C41) % 10000;
    boolean flag = false;
    string label = "rec745";
};
type R746 record {
    int id;
    int a = C42 + 746;
    int b = C63 * 5;
    int c = (746 + C42) % 10000;
    boolean flag = true;
    string label = "rec746";
};
type R747 record {
    int id;
    int a = C43 + 747;
    int b = C2 * 6;
    int c = (747 + C43) % 10000;
    boolean flag = false;
    string label = "rec747";
};
type R748 record {
    int id;
    int a = C44 + 748;
    int b = C5 * 7;
    int c = (748 + C44) % 10000;
    boolean flag = true;
    string label = "rec748";
};
type R749 record {
    int id;
    int a = C45 + 749;
    int b = C8 * 1;
    int c = (749 + C45) % 10000;
    boolean flag = false;
    string label = "rec749";
};
type R750 record {
    int id;
    int a = C46 + 750;
    int b = C11 * 2;
    int c = (750 + C46) % 10000;
    boolean flag = true;
    string label = "rec750";
};
type R751 record {
    int id;
    int a = C47 + 751;
    int b = C14 * 3;
    int c = (751 + C47) % 10000;
    boolean flag = false;
    string label = "rec751";
};
type R752 record {
    int id;
    int a = C48 + 752;
    int b = C17 * 4;
    int c = (752 + C48) % 10000;
    boolean flag = true;
    string label = "rec752";
};
type R753 record {
    int id;
    int a = C49 + 753;
    int b = C20 * 5;
    int c = (753 + C49) % 10000;
    boolean flag = false;
    string label = "rec753";
};
type R754 record {
    int id;
    int a = C50 + 754;
    int b = C23 * 6;
    int c = (754 + C50) % 10000;
    boolean flag = true;
    string label = "rec754";
};
type R755 record {
    int id;
    int a = C51 + 755;
    int b = C26 * 7;
    int c = (755 + C51) % 10000;
    boolean flag = false;
    string label = "rec755";
};
type R756 record {
    int id;
    int a = C52 + 756;
    int b = C29 * 1;
    int c = (756 + C52) % 10000;
    boolean flag = true;
    string label = "rec756";
};
type R757 record {
    int id;
    int a = C53 + 757;
    int b = C32 * 2;
    int c = (757 + C53) % 10000;
    boolean flag = false;
    string label = "rec757";
};
type R758 record {
    int id;
    int a = C54 + 758;
    int b = C35 * 3;
    int c = (758 + C54) % 10000;
    boolean flag = true;
    string label = "rec758";
};
type R759 record {
    int id;
    int a = C55 + 759;
    int b = C38 * 4;
    int c = (759 + C55) % 10000;
    boolean flag = false;
    string label = "rec759";
};
type R760 record {
    int id;
    int a = C56 + 760;
    int b = C41 * 5;
    int c = (760 + C56) % 10000;
    boolean flag = true;
    string label = "rec760";
};
type R761 record {
    int id;
    int a = C57 + 761;
    int b = C44 * 6;
    int c = (761 + C57) % 10000;
    boolean flag = false;
    string label = "rec761";
};
type R762 record {
    int id;
    int a = C58 + 762;
    int b = C47 * 7;
    int c = (762 + C58) % 10000;
    boolean flag = true;
    string label = "rec762";
};
type R763 record {
    int id;
    int a = C59 + 763;
    int b = C50 * 1;
    int c = (763 + C59) % 10000;
    boolean flag = false;
    string label = "rec763";
};
type R764 record {
    int id;
    int a = C60 + 764;
    int b = C53 * 2;
    int c = (764 + C60) % 10000;
    boolean flag = true;
    string label = "rec764";
};
type R765 record {
    int id;
    int a = C61 + 765;
    int b = C56 * 3;
    int c = (765 + C61) % 10000;
    boolean flag = false;
    string label = "rec765";
};
type R766 record {
    int id;
    int a = C62 + 766;
    int b = C59 * 4;
    int c = (766 + C62) % 10000;
    boolean flag = true;
    string label = "rec766";
};
type R767 record {
    int id;
    int a = C63 + 767;
    int b = C62 * 5;
    int c = (767 + C63) % 10000;
    boolean flag = false;
    string label = "rec767";
};
type R768 record {
    int id;
    int a = C0 + 768;
    int b = C1 * 6;
    int c = (768 + C0) % 10000;
    boolean flag = true;
    string label = "rec768";
};
type R769 record {
    int id;
    int a = C1 + 769;
    int b = C4 * 7;
    int c = (769 + C1) % 10000;
    boolean flag = false;
    string label = "rec769";
};
type R770 record {
    int id;
    int a = C2 + 770;
    int b = C7 * 1;
    int c = (770 + C2) % 10000;
    boolean flag = true;
    string label = "rec770";
};
type R771 record {
    int id;
    int a = C3 + 771;
    int b = C10 * 2;
    int c = (771 + C3) % 10000;
    boolean flag = false;
    string label = "rec771";
};
type R772 record {
    int id;
    int a = C4 + 772;
    int b = C13 * 3;
    int c = (772 + C4) % 10000;
    boolean flag = true;
    string label = "rec772";
};
type R773 record {
    int id;
    int a = C5 + 773;
    int b = C16 * 4;
    int c = (773 + C5) % 10000;
    boolean flag = false;
    string label = "rec773";
};
type R774 record {
    int id;
    int a = C6 + 774;
    int b = C19 * 5;
    int c = (774 + C6) % 10000;
    boolean flag = true;
    string label = "rec774";
};
type R775 record {
    int id;
    int a = C7 + 775;
    int b = C22 * 6;
    int c = (775 + C7) % 10000;
    boolean flag = false;
    string label = "rec775";
};
type R776 record {
    int id;
    int a = C8 + 776;
    int b = C25 * 7;
    int c = (776 + C8) % 10000;
    boolean flag = true;
    string label = "rec776";
};
type R777 record {
    int id;
    int a = C9 + 777;
    int b = C28 * 1;
    int c = (777 + C9) % 10000;
    boolean flag = false;
    string label = "rec777";
};
type R778 record {
    int id;
    int a = C10 + 778;
    int b = C31 * 2;
    int c = (778 + C10) % 10000;
    boolean flag = true;
    string label = "rec778";
};
type R779 record {
    int id;
    int a = C11 + 779;
    int b = C34 * 3;
    int c = (779 + C11) % 10000;
    boolean flag = false;
    string label = "rec779";
};
type R780 record {
    int id;
    int a = C12 + 780;
    int b = C37 * 4;
    int c = (780 + C12) % 10000;
    boolean flag = true;
    string label = "rec780";
};
type R781 record {
    int id;
    int a = C13 + 781;
    int b = C40 * 5;
    int c = (781 + C13) % 10000;
    boolean flag = false;
    string label = "rec781";
};
type R782 record {
    int id;
    int a = C14 + 782;
    int b = C43 * 6;
    int c = (782 + C14) % 10000;
    boolean flag = true;
    string label = "rec782";
};
type R783 record {
    int id;
    int a = C15 + 783;
    int b = C46 * 7;
    int c = (783 + C15) % 10000;
    boolean flag = false;
    string label = "rec783";
};
type R784 record {
    int id;
    int a = C16 + 784;
    int b = C49 * 1;
    int c = (784 + C16) % 10000;
    boolean flag = true;
    string label = "rec784";
};
type R785 record {
    int id;
    int a = C17 + 785;
    int b = C52 * 2;
    int c = (785 + C17) % 10000;
    boolean flag = false;
    string label = "rec785";
};
type R786 record {
    int id;
    int a = C18 + 786;
    int b = C55 * 3;
    int c = (786 + C18) % 10000;
    boolean flag = true;
    string label = "rec786";
};
type R787 record {
    int id;
    int a = C19 + 787;
    int b = C58 * 4;
    int c = (787 + C19) % 10000;
    boolean flag = false;
    string label = "rec787";
};
type R788 record {
    int id;
    int a = C20 + 788;
    int b = C61 * 5;
    int c = (788 + C20) % 10000;
    boolean flag = true;
    string label = "rec788";
};
type R789 record {
    int id;
    int a = C21 + 789;
    int b = C0 * 6;
    int c = (789 + C21) % 10000;
    boolean flag = false;
    string label = "rec789";
};
type R790 record {
    int id;
    int a = C22 + 790;
    int b = C3 * 7;
    int c = (790 + C22) % 10000;
    boolean flag = true;
    string label = "rec790";
};
type R791 record {
    int id;
    int a = C23 + 791;
    int b = C6 * 1;
    int c = (791 + C23) % 10000;
    boolean flag = false;
    string label = "rec791";
};
type R792 record {
    int id;
    int a = C24 + 792;
    int b = C9 * 2;
    int c = (792 + C24) % 10000;
    boolean flag = true;
    string label = "rec792";
};
type R793 record {
    int id;
    int a = C25 + 793;
    int b = C12 * 3;
    int c = (793 + C25) % 10000;
    boolean flag = false;
    string label = "rec793";
};
type R794 record {
    int id;
    int a = C26 + 794;
    int b = C15 * 4;
    int c = (794 + C26) % 10000;
    boolean flag = true;
    string label = "rec794";
};
type R795 record {
    int id;
    int a = C27 + 795;
    int b = C18 * 5;
    int c = (795 + C27) % 10000;
    boolean flag = false;
    string label = "rec795";
};
type R796 record {
    int id;
    int a = C28 + 796;
    int b = C21 * 6;
    int c = (796 + C28) % 10000;
    boolean flag = true;
    string label = "rec796";
};
type R797 record {
    int id;
    int a = C29 + 797;
    int b = C24 * 7;
    int c = (797 + C29) % 10000;
    boolean flag = false;
    string label = "rec797";
};
type R798 record {
    int id;
    int a = C30 + 798;
    int b = C27 * 1;
    int c = (798 + C30) % 10000;
    boolean flag = true;
    string label = "rec798";
};
type R799 record {
    int id;
    int a = C31 + 799;
    int b = C30 * 2;
    int c = (799 + C31) % 10000;
    boolean flag = false;
    string label = "rec799";
};
type R800 record {
    int id;
    int a = C32 + 800;
    int b = C33 * 3;
    int c = (800 + C32) % 10000;
    boolean flag = true;
    string label = "rec800";
};
type R801 record {
    int id;
    int a = C33 + 801;
    int b = C36 * 4;
    int c = (801 + C33) % 10000;
    boolean flag = false;
    string label = "rec801";
};
type R802 record {
    int id;
    int a = C34 + 802;
    int b = C39 * 5;
    int c = (802 + C34) % 10000;
    boolean flag = true;
    string label = "rec802";
};
type R803 record {
    int id;
    int a = C35 + 803;
    int b = C42 * 6;
    int c = (803 + C35) % 10000;
    boolean flag = false;
    string label = "rec803";
};
type R804 record {
    int id;
    int a = C36 + 804;
    int b = C45 * 7;
    int c = (804 + C36) % 10000;
    boolean flag = true;
    string label = "rec804";
};
type R805 record {
    int id;
    int a = C37 + 805;
    int b = C48 * 1;
    int c = (805 + C37) % 10000;
    boolean flag = false;
    string label = "rec805";
};
type R806 record {
    int id;
    int a = C38 + 806;
    int b = C51 * 2;
    int c = (806 + C38) % 10000;
    boolean flag = true;
    string label = "rec806";
};
type R807 record {
    int id;
    int a = C39 + 807;
    int b = C54 * 3;
    int c = (807 + C39) % 10000;
    boolean flag = false;
    string label = "rec807";
};
type R808 record {
    int id;
    int a = C40 + 808;
    int b = C57 * 4;
    int c = (808 + C40) % 10000;
    boolean flag = true;
    string label = "rec808";
};
type R809 record {
    int id;
    int a = C41 + 809;
    int b = C60 * 5;
    int c = (809 + C41) % 10000;
    boolean flag = false;
    string label = "rec809";
};
type R810 record {
    int id;
    int a = C42 + 810;
    int b = C63 * 6;
    int c = (810 + C42) % 10000;
    boolean flag = true;
    string label = "rec810";
};
type R811 record {
    int id;
    int a = C43 + 811;
    int b = C2 * 7;
    int c = (811 + C43) % 10000;
    boolean flag = false;
    string label = "rec811";
};
type R812 record {
    int id;
    int a = C44 + 812;
    int b = C5 * 1;
    int c = (812 + C44) % 10000;
    boolean flag = true;
    string label = "rec812";
};
type R813 record {
    int id;
    int a = C45 + 813;
    int b = C8 * 2;
    int c = (813 + C45) % 10000;
    boolean flag = false;
    string label = "rec813";
};
type R814 record {
    int id;
    int a = C46 + 814;
    int b = C11 * 3;
    int c = (814 + C46) % 10000;
    boolean flag = true;
    string label = "rec814";
};
type R815 record {
    int id;
    int a = C47 + 815;
    int b = C14 * 4;
    int c = (815 + C47) % 10000;
    boolean flag = false;
    string label = "rec815";
};
type R816 record {
    int id;
    int a = C48 + 816;
    int b = C17 * 5;
    int c = (816 + C48) % 10000;
    boolean flag = true;
    string label = "rec816";
};
type R817 record {
    int id;
    int a = C49 + 817;
    int b = C20 * 6;
    int c = (817 + C49) % 10000;
    boolean flag = false;
    string label = "rec817";
};
type R818 record {
    int id;
    int a = C50 + 818;
    int b = C23 * 7;
    int c = (818 + C50) % 10000;
    boolean flag = true;
    string label = "rec818";
};
type R819 record {
    int id;
    int a = C51 + 819;
    int b = C26 * 1;
    int c = (819 + C51) % 10000;
    boolean flag = false;
    string label = "rec819";
};
type R820 record {
    int id;
    int a = C52 + 820;
    int b = C29 * 2;
    int c = (820 + C52) % 10000;
    boolean flag = true;
    string label = "rec820";
};
type R821 record {
    int id;
    int a = C53 + 821;
    int b = C32 * 3;
    int c = (821 + C53) % 10000;
    boolean flag = false;
    string label = "rec821";
};
type R822 record {
    int id;
    int a = C54 + 822;
    int b = C35 * 4;
    int c = (822 + C54) % 10000;
    boolean flag = true;
    string label = "rec822";
};
type R823 record {
    int id;
    int a = C55 + 823;
    int b = C38 * 5;
    int c = (823 + C55) % 10000;
    boolean flag = false;
    string label = "rec823";
};
type R824 record {
    int id;
    int a = C56 + 824;
    int b = C41 * 6;
    int c = (824 + C56) % 10000;
    boolean flag = true;
    string label = "rec824";
};
type R825 record {
    int id;
    int a = C57 + 825;
    int b = C44 * 7;
    int c = (825 + C57) % 10000;
    boolean flag = false;
    string label = "rec825";
};
type R826 record {
    int id;
    int a = C58 + 826;
    int b = C47 * 1;
    int c = (826 + C58) % 10000;
    boolean flag = true;
    string label = "rec826";
};
type R827 record {
    int id;
    int a = C59 + 827;
    int b = C50 * 2;
    int c = (827 + C59) % 10000;
    boolean flag = false;
    string label = "rec827";
};
type R828 record {
    int id;
    int a = C60 + 828;
    int b = C53 * 3;
    int c = (828 + C60) % 10000;
    boolean flag = true;
    string label = "rec828";
};
type R829 record {
    int id;
    int a = C61 + 829;
    int b = C56 * 4;
    int c = (829 + C61) % 10000;
    boolean flag = false;
    string label = "rec829";
};
type R830 record {
    int id;
    int a = C62 + 830;
    int b = C59 * 5;
    int c = (830 + C62) % 10000;
    boolean flag = true;
    string label = "rec830";
};
type R831 record {
    int id;
    int a = C63 + 831;
    int b = C62 * 6;
    int c = (831 + C63) % 10000;
    boolean flag = false;
    string label = "rec831";
};
type R832 record {
    int id;
    int a = C0 + 832;
    int b = C1 * 7;
    int c = (832 + C0) % 10000;
    boolean flag = true;
    string label = "rec832";
};
type R833 record {
    int id;
    int a = C1 + 833;
    int b = C4 * 1;
    int c = (833 + C1) % 10000;
    boolean flag = false;
    string label = "rec833";
};
type R834 record {
    int id;
    int a = C2 + 834;
    int b = C7 * 2;
    int c = (834 + C2) % 10000;
    boolean flag = true;
    string label = "rec834";
};
type R835 record {
    int id;
    int a = C3 + 835;
    int b = C10 * 3;
    int c = (835 + C3) % 10000;
    boolean flag = false;
    string label = "rec835";
};
type R836 record {
    int id;
    int a = C4 + 836;
    int b = C13 * 4;
    int c = (836 + C4) % 10000;
    boolean flag = true;
    string label = "rec836";
};
type R837 record {
    int id;
    int a = C5 + 837;
    int b = C16 * 5;
    int c = (837 + C5) % 10000;
    boolean flag = false;
    string label = "rec837";
};
type R838 record {
    int id;
    int a = C6 + 838;
    int b = C19 * 6;
    int c = (838 + C6) % 10000;
    boolean flag = true;
    string label = "rec838";
};
type R839 record {
    int id;
    int a = C7 + 839;
    int b = C22 * 7;
    int c = (839 + C7) % 10000;
    boolean flag = false;
    string label = "rec839";
};
type R840 record {
    int id;
    int a = C8 + 840;
    int b = C25 * 1;
    int c = (840 + C8) % 10000;
    boolean flag = true;
    string label = "rec840";
};
type R841 record {
    int id;
    int a = C9 + 841;
    int b = C28 * 2;
    int c = (841 + C9) % 10000;
    boolean flag = false;
    string label = "rec841";
};
type R842 record {
    int id;
    int a = C10 + 842;
    int b = C31 * 3;
    int c = (842 + C10) % 10000;
    boolean flag = true;
    string label = "rec842";
};
type R843 record {
    int id;
    int a = C11 + 843;
    int b = C34 * 4;
    int c = (843 + C11) % 10000;
    boolean flag = false;
    string label = "rec843";
};
type R844 record {
    int id;
    int a = C12 + 844;
    int b = C37 * 5;
    int c = (844 + C12) % 10000;
    boolean flag = true;
    string label = "rec844";
};
type R845 record {
    int id;
    int a = C13 + 845;
    int b = C40 * 6;
    int c = (845 + C13) % 10000;
    boolean flag = false;
    string label = "rec845";
};
type R846 record {
    int id;
    int a = C14 + 846;
    int b = C43 * 7;
    int c = (846 + C14) % 10000;
    boolean flag = true;
    string label = "rec846";
};
type R847 record {
    int id;
    int a = C15 + 847;
    int b = C46 * 1;
    int c = (847 + C15) % 10000;
    boolean flag = false;
    string label = "rec847";
};
type R848 record {
    int id;
    int a = C16 + 848;
    int b = C49 * 2;
    int c = (848 + C16) % 10000;
    boolean flag = true;
    string label = "rec848";
};
type R849 record {
    int id;
    int a = C17 + 849;
    int b = C52 * 3;
    int c = (849 + C17) % 10000;
    boolean flag = false;
    string label = "rec849";
};
type R850 record {
    int id;
    int a = C18 + 850;
    int b = C55 * 4;
    int c = (850 + C18) % 10000;
    boolean flag = true;
    string label = "rec850";
};
type R851 record {
    int id;
    int a = C19 + 851;
    int b = C58 * 5;
    int c = (851 + C19) % 10000;
    boolean flag = false;
    string label = "rec851";
};
type R852 record {
    int id;
    int a = C20 + 852;
    int b = C61 * 6;
    int c = (852 + C20) % 10000;
    boolean flag = true;
    string label = "rec852";
};
type R853 record {
    int id;
    int a = C21 + 853;
    int b = C0 * 7;
    int c = (853 + C21) % 10000;
    boolean flag = false;
    string label = "rec853";
};
type R854 record {
    int id;
    int a = C22 + 854;
    int b = C3 * 1;
    int c = (854 + C22) % 10000;
    boolean flag = true;
    string label = "rec854";
};
type R855 record {
    int id;
    int a = C23 + 855;
    int b = C6 * 2;
    int c = (855 + C23) % 10000;
    boolean flag = false;
    string label = "rec855";
};
type R856 record {
    int id;
    int a = C24 + 856;
    int b = C9 * 3;
    int c = (856 + C24) % 10000;
    boolean flag = true;
    string label = "rec856";
};
type R857 record {
    int id;
    int a = C25 + 857;
    int b = C12 * 4;
    int c = (857 + C25) % 10000;
    boolean flag = false;
    string label = "rec857";
};
type R858 record {
    int id;
    int a = C26 + 858;
    int b = C15 * 5;
    int c = (858 + C26) % 10000;
    boolean flag = true;
    string label = "rec858";
};
type R859 record {
    int id;
    int a = C27 + 859;
    int b = C18 * 6;
    int c = (859 + C27) % 10000;
    boolean flag = false;
    string label = "rec859";
};
type R860 record {
    int id;
    int a = C28 + 860;
    int b = C21 * 7;
    int c = (860 + C28) % 10000;
    boolean flag = true;
    string label = "rec860";
};
type R861 record {
    int id;
    int a = C29 + 861;
    int b = C24 * 1;
    int c = (861 + C29) % 10000;
    boolean flag = false;
    string label = "rec861";
};
type R862 record {
    int id;
    int a = C30 + 862;
    int b = C27 * 2;
    int c = (862 + C30) % 10000;
    boolean flag = true;
    string label = "rec862";
};
type R863 record {
    int id;
    int a = C31 + 863;
    int b = C30 * 3;
    int c = (863 + C31) % 10000;
    boolean flag = false;
    string label = "rec863";
};
type R864 record {
    int id;
    int a = C32 + 864;
    int b = C33 * 4;
    int c = (864 + C32) % 10000;
    boolean flag = true;
    string label = "rec864";
};
type R865 record {
    int id;
    int a = C33 + 865;
    int b = C36 * 5;
    int c = (865 + C33) % 10000;
    boolean flag = false;
    string label = "rec865";
};
type R866 record {
    int id;
    int a = C34 + 866;
    int b = C39 * 6;
    int c = (866 + C34) % 10000;
    boolean flag = true;
    string label = "rec866";
};
type R867 record {
    int id;
    int a = C35 + 867;
    int b = C42 * 7;
    int c = (867 + C35) % 10000;
    boolean flag = false;
    string label = "rec867";
};
type R868 record {
    int id;
    int a = C36 + 868;
    int b = C45 * 1;
    int c = (868 + C36) % 10000;
    boolean flag = true;
    string label = "rec868";
};
type R869 record {
    int id;
    int a = C37 + 869;
    int b = C48 * 2;
    int c = (869 + C37) % 10000;
    boolean flag = false;
    string label = "rec869";
};
type R870 record {
    int id;
    int a = C38 + 870;
    int b = C51 * 3;
    int c = (870 + C38) % 10000;
    boolean flag = true;
    string label = "rec870";
};
type R871 record {
    int id;
    int a = C39 + 871;
    int b = C54 * 4;
    int c = (871 + C39) % 10000;
    boolean flag = false;
    string label = "rec871";
};
type R872 record {
    int id;
    int a = C40 + 872;
    int b = C57 * 5;
    int c = (872 + C40) % 10000;
    boolean flag = true;
    string label = "rec872";
};
type R873 record {
    int id;
    int a = C41 + 873;
    int b = C60 * 6;
    int c = (873 + C41) % 10000;
    boolean flag = false;
    string label = "rec873";
};
type R874 record {
    int id;
    int a = C42 + 874;
    int b = C63 * 7;
    int c = (874 + C42) % 10000;
    boolean flag = true;
    string label = "rec874";
};
type R875 record {
    int id;
    int a = C43 + 875;
    int b = C2 * 1;
    int c = (875 + C43) % 10000;
    boolean flag = false;
    string label = "rec875";
};
type R876 record {
    int id;
    int a = C44 + 876;
    int b = C5 * 2;
    int c = (876 + C44) % 10000;
    boolean flag = true;
    string label = "rec876";
};
type R877 record {
    int id;
    int a = C45 + 877;
    int b = C8 * 3;
    int c = (877 + C45) % 10000;
    boolean flag = false;
    string label = "rec877";
};
type R878 record {
    int id;
    int a = C46 + 878;
    int b = C11 * 4;
    int c = (878 + C46) % 10000;
    boolean flag = true;
    string label = "rec878";
};
type R879 record {
    int id;
    int a = C47 + 879;
    int b = C14 * 5;
    int c = (879 + C47) % 10000;
    boolean flag = false;
    string label = "rec879";
};
type R880 record {
    int id;
    int a = C48 + 880;
    int b = C17 * 6;
    int c = (880 + C48) % 10000;
    boolean flag = true;
    string label = "rec880";
};
type R881 record {
    int id;
    int a = C49 + 881;
    int b = C20 * 7;
    int c = (881 + C49) % 10000;
    boolean flag = false;
    string label = "rec881";
};
type R882 record {
    int id;
    int a = C50 + 882;
    int b = C23 * 1;
    int c = (882 + C50) % 10000;
    boolean flag = true;
    string label = "rec882";
};
type R883 record {
    int id;
    int a = C51 + 883;
    int b = C26 * 2;
    int c = (883 + C51) % 10000;
    boolean flag = false;
    string label = "rec883";
};
type R884 record {
    int id;
    int a = C52 + 884;
    int b = C29 * 3;
    int c = (884 + C52) % 10000;
    boolean flag = true;
    string label = "rec884";
};
type R885 record {
    int id;
    int a = C53 + 885;
    int b = C32 * 4;
    int c = (885 + C53) % 10000;
    boolean flag = false;
    string label = "rec885";
};
type R886 record {
    int id;
    int a = C54 + 886;
    int b = C35 * 5;
    int c = (886 + C54) % 10000;
    boolean flag = true;
    string label = "rec886";
};
type R887 record {
    int id;
    int a = C55 + 887;
    int b = C38 * 6;
    int c = (887 + C55) % 10000;
    boolean flag = false;
    string label = "rec887";
};
type R888 record {
    int id;
    int a = C56 + 888;
    int b = C41 * 7;
    int c = (888 + C56) % 10000;
    boolean flag = true;
    string label = "rec888";
};
type R889 record {
    int id;
    int a = C57 + 889;
    int b = C44 * 1;
    int c = (889 + C57) % 10000;
    boolean flag = false;
    string label = "rec889";
};
type R890 record {
    int id;
    int a = C58 + 890;
    int b = C47 * 2;
    int c = (890 + C58) % 10000;
    boolean flag = true;
    string label = "rec890";
};
type R891 record {
    int id;
    int a = C59 + 891;
    int b = C50 * 3;
    int c = (891 + C59) % 10000;
    boolean flag = false;
    string label = "rec891";
};
type R892 record {
    int id;
    int a = C60 + 892;
    int b = C53 * 4;
    int c = (892 + C60) % 10000;
    boolean flag = true;
    string label = "rec892";
};
type R893 record {
    int id;
    int a = C61 + 893;
    int b = C56 * 5;
    int c = (893 + C61) % 10000;
    boolean flag = false;
    string label = "rec893";
};
type R894 record {
    int id;
    int a = C62 + 894;
    int b = C59 * 6;
    int c = (894 + C62) % 10000;
    boolean flag = true;
    string label = "rec894";
};
type R895 record {
    int id;
    int a = C63 + 895;
    int b = C62 * 7;
    int c = (895 + C63) % 10000;
    boolean flag = false;
    string label = "rec895";
};
type R896 record {
    int id;
    int a = C0 + 896;
    int b = C1 * 1;
    int c = (896 + C0) % 10000;
    boolean flag = true;
    string label = "rec896";
};
type R897 record {
    int id;
    int a = C1 + 897;
    int b = C4 * 2;
    int c = (897 + C1) % 10000;
    boolean flag = false;
    string label = "rec897";
};
type R898 record {
    int id;
    int a = C2 + 898;
    int b = C7 * 3;
    int c = (898 + C2) % 10000;
    boolean flag = true;
    string label = "rec898";
};
type R899 record {
    int id;
    int a = C3 + 899;
    int b = C10 * 4;
    int c = (899 + C3) % 10000;
    boolean flag = false;
    string label = "rec899";
};
type R900 record {
    int id;
    int a = C4 + 900;
    int b = C13 * 5;
    int c = (900 + C4) % 10000;
    boolean flag = true;
    string label = "rec900";
};
type R901 record {
    int id;
    int a = C5 + 901;
    int b = C16 * 6;
    int c = (901 + C5) % 10000;
    boolean flag = false;
    string label = "rec901";
};
type R902 record {
    int id;
    int a = C6 + 902;
    int b = C19 * 7;
    int c = (902 + C6) % 10000;
    boolean flag = true;
    string label = "rec902";
};
type R903 record {
    int id;
    int a = C7 + 903;
    int b = C22 * 1;
    int c = (903 + C7) % 10000;
    boolean flag = false;
    string label = "rec903";
};
type R904 record {
    int id;
    int a = C8 + 904;
    int b = C25 * 2;
    int c = (904 + C8) % 10000;
    boolean flag = true;
    string label = "rec904";
};
type R905 record {
    int id;
    int a = C9 + 905;
    int b = C28 * 3;
    int c = (905 + C9) % 10000;
    boolean flag = false;
    string label = "rec905";
};
type R906 record {
    int id;
    int a = C10 + 906;
    int b = C31 * 4;
    int c = (906 + C10) % 10000;
    boolean flag = true;
    string label = "rec906";
};
type R907 record {
    int id;
    int a = C11 + 907;
    int b = C34 * 5;
    int c = (907 + C11) % 10000;
    boolean flag = false;
    string label = "rec907";
};
type R908 record {
    int id;
    int a = C12 + 908;
    int b = C37 * 6;
    int c = (908 + C12) % 10000;
    boolean flag = true;
    string label = "rec908";
};
type R909 record {
    int id;
    int a = C13 + 909;
    int b = C40 * 7;
    int c = (909 + C13) % 10000;
    boolean flag = false;
    string label = "rec909";
};
type R910 record {
    int id;
    int a = C14 + 910;
    int b = C43 * 1;
    int c = (910 + C14) % 10000;
    boolean flag = true;
    string label = "rec910";
};
type R911 record {
    int id;
    int a = C15 + 911;
    int b = C46 * 2;
    int c = (911 + C15) % 10000;
    boolean flag = false;
    string label = "rec911";
};
type R912 record {
    int id;
    int a = C16 + 912;
    int b = C49 * 3;
    int c = (912 + C16) % 10000;
    boolean flag = true;
    string label = "rec912";
};
type R913 record {
    int id;
    int a = C17 + 913;
    int b = C52 * 4;
    int c = (913 + C17) % 10000;
    boolean flag = false;
    string label = "rec913";
};
type R914 record {
    int id;
    int a = C18 + 914;
    int b = C55 * 5;
    int c = (914 + C18) % 10000;
    boolean flag = true;
    string label = "rec914";
};
type R915 record {
    int id;
    int a = C19 + 915;
    int b = C58 * 6;
    int c = (915 + C19) % 10000;
    boolean flag = false;
    string label = "rec915";
};
type R916 record {
    int id;
    int a = C20 + 916;
    int b = C61 * 7;
    int c = (916 + C20) % 10000;
    boolean flag = true;
    string label = "rec916";
};
type R917 record {
    int id;
    int a = C21 + 917;
    int b = C0 * 1;
    int c = (917 + C21) % 10000;
    boolean flag = false;
    string label = "rec917";
};
type R918 record {
    int id;
    int a = C22 + 918;
    int b = C3 * 2;
    int c = (918 + C22) % 10000;
    boolean flag = true;
    string label = "rec918";
};
type R919 record {
    int id;
    int a = C23 + 919;
    int b = C6 * 3;
    int c = (919 + C23) % 10000;
    boolean flag = false;
    string label = "rec919";
};
type R920 record {
    int id;
    int a = C24 + 920;
    int b = C9 * 4;
    int c = (920 + C24) % 10000;
    boolean flag = true;
    string label = "rec920";
};
type R921 record {
    int id;
    int a = C25 + 921;
    int b = C12 * 5;
    int c = (921 + C25) % 10000;
    boolean flag = false;
    string label = "rec921";
};
type R922 record {
    int id;
    int a = C26 + 922;
    int b = C15 * 6;
    int c = (922 + C26) % 10000;
    boolean flag = true;
    string label = "rec922";
};
type R923 record {
    int id;
    int a = C27 + 923;
    int b = C18 * 7;
    int c = (923 + C27) % 10000;
    boolean flag = false;
    string label = "rec923";
};
type R924 record {
    int id;
    int a = C28 + 924;
    int b = C21 * 1;
    int c = (924 + C28) % 10000;
    boolean flag = true;
    string label = "rec924";
};
type R925 record {
    int id;
    int a = C29 + 925;
    int b = C24 * 2;
    int c = (925 + C29) % 10000;
    boolean flag = false;
    string label = "rec925";
};
type R926 record {
    int id;
    int a = C30 + 926;
    int b = C27 * 3;
    int c = (926 + C30) % 10000;
    boolean flag = true;
    string label = "rec926";
};
type R927 record {
    int id;
    int a = C31 + 927;
    int b = C30 * 4;
    int c = (927 + C31) % 10000;
    boolean flag = false;
    string label = "rec927";
};
type R928 record {
    int id;
    int a = C32 + 928;
    int b = C33 * 5;
    int c = (928 + C32) % 10000;
    boolean flag = true;
    string label = "rec928";
};
type R929 record {
    int id;
    int a = C33 + 929;
    int b = C36 * 6;
    int c = (929 + C33) % 10000;
    boolean flag = false;
    string label = "rec929";
};
type R930 record {
    int id;
    int a = C34 + 930;
    int b = C39 * 7;
    int c = (930 + C34) % 10000;
    boolean flag = true;
    string label = "rec930";
};
type R931 record {
    int id;
    int a = C35 + 931;
    int b = C42 * 1;
    int c = (931 + C35) % 10000;
    boolean flag = false;
    string label = "rec931";
};
type R932 record {
    int id;
    int a = C36 + 932;
    int b = C45 * 2;
    int c = (932 + C36) % 10000;
    boolean flag = true;
    string label = "rec932";
};
type R933 record {
    int id;
    int a = C37 + 933;
    int b = C48 * 3;
    int c = (933 + C37) % 10000;
    boolean flag = false;
    string label = "rec933";
};
type R934 record {
    int id;
    int a = C38 + 934;
    int b = C51 * 4;
    int c = (934 + C38) % 10000;
    boolean flag = true;
    string label = "rec934";
};
type R935 record {
    int id;
    int a = C39 + 935;
    int b = C54 * 5;
    int c = (935 + C39) % 10000;
    boolean flag = false;
    string label = "rec935";
};
type R936 record {
    int id;
    int a = C40 + 936;
    int b = C57 * 6;
    int c = (936 + C40) % 10000;
    boolean flag = true;
    string label = "rec936";
};
type R937 record {
    int id;
    int a = C41 + 937;
    int b = C60 * 7;
    int c = (937 + C41) % 10000;
    boolean flag = false;
    string label = "rec937";
};
type R938 record {
    int id;
    int a = C42 + 938;
    int b = C63 * 1;
    int c = (938 + C42) % 10000;
    boolean flag = true;
    string label = "rec938";
};
type R939 record {
    int id;
    int a = C43 + 939;
    int b = C2 * 2;
    int c = (939 + C43) % 10000;
    boolean flag = false;
    string label = "rec939";
};
type R940 record {
    int id;
    int a = C44 + 940;
    int b = C5 * 3;
    int c = (940 + C44) % 10000;
    boolean flag = true;
    string label = "rec940";
};
type R941 record {
    int id;
    int a = C45 + 941;
    int b = C8 * 4;
    int c = (941 + C45) % 10000;
    boolean flag = false;
    string label = "rec941";
};
type R942 record {
    int id;
    int a = C46 + 942;
    int b = C11 * 5;
    int c = (942 + C46) % 10000;
    boolean flag = true;
    string label = "rec942";
};
type R943 record {
    int id;
    int a = C47 + 943;
    int b = C14 * 6;
    int c = (943 + C47) % 10000;
    boolean flag = false;
    string label = "rec943";
};
type R944 record {
    int id;
    int a = C48 + 944;
    int b = C17 * 7;
    int c = (944 + C48) % 10000;
    boolean flag = true;
    string label = "rec944";
};
type R945 record {
    int id;
    int a = C49 + 945;
    int b = C20 * 1;
    int c = (945 + C49) % 10000;
    boolean flag = false;
    string label = "rec945";
};
type R946 record {
    int id;
    int a = C50 + 946;
    int b = C23 * 2;
    int c = (946 + C50) % 10000;
    boolean flag = true;
    string label = "rec946";
};
type R947 record {
    int id;
    int a = C51 + 947;
    int b = C26 * 3;
    int c = (947 + C51) % 10000;
    boolean flag = false;
    string label = "rec947";
};
type R948 record {
    int id;
    int a = C52 + 948;
    int b = C29 * 4;
    int c = (948 + C52) % 10000;
    boolean flag = true;
    string label = "rec948";
};
type R949 record {
    int id;
    int a = C53 + 949;
    int b = C32 * 5;
    int c = (949 + C53) % 10000;
    boolean flag = false;
    string label = "rec949";
};
type R950 record {
    int id;
    int a = C54 + 950;
    int b = C35 * 6;
    int c = (950 + C54) % 10000;
    boolean flag = true;
    string label = "rec950";
};
type R951 record {
    int id;
    int a = C55 + 951;
    int b = C38 * 7;
    int c = (951 + C55) % 10000;
    boolean flag = false;
    string label = "rec951";
};
type R952 record {
    int id;
    int a = C56 + 952;
    int b = C41 * 1;
    int c = (952 + C56) % 10000;
    boolean flag = true;
    string label = "rec952";
};
type R953 record {
    int id;
    int a = C57 + 953;
    int b = C44 * 2;
    int c = (953 + C57) % 10000;
    boolean flag = false;
    string label = "rec953";
};
type R954 record {
    int id;
    int a = C58 + 954;
    int b = C47 * 3;
    int c = (954 + C58) % 10000;
    boolean flag = true;
    string label = "rec954";
};
type R955 record {
    int id;
    int a = C59 + 955;
    int b = C50 * 4;
    int c = (955 + C59) % 10000;
    boolean flag = false;
    string label = "rec955";
};
type R956 record {
    int id;
    int a = C60 + 956;
    int b = C53 * 5;
    int c = (956 + C60) % 10000;
    boolean flag = true;
    string label = "rec956";
};
type R957 record {
    int id;
    int a = C61 + 957;
    int b = C56 * 6;
    int c = (957 + C61) % 10000;
    boolean flag = false;
    string label = "rec957";
};
type R958 record {
    int id;
    int a = C62 + 958;
    int b = C59 * 7;
    int c = (958 + C62) % 10000;
    boolean flag = true;
    string label = "rec958";
};
type R959 record {
    int id;
    int a = C63 + 959;
    int b = C62 * 1;
    int c = (959 + C63) % 10000;
    boolean flag = false;
    string label = "rec959";
};
type R960 record {
    int id;
    int a = C0 + 960;
    int b = C1 * 2;
    int c = (960 + C0) % 10000;
    boolean flag = true;
    string label = "rec960";
};
type R961 record {
    int id;
    int a = C1 + 961;
    int b = C4 * 3;
    int c = (961 + C1) % 10000;
    boolean flag = false;
    string label = "rec961";
};
type R962 record {
    int id;
    int a = C2 + 962;
    int b = C7 * 4;
    int c = (962 + C2) % 10000;
    boolean flag = true;
    string label = "rec962";
};
type R963 record {
    int id;
    int a = C3 + 963;
    int b = C10 * 5;
    int c = (963 + C3) % 10000;
    boolean flag = false;
    string label = "rec963";
};
type R964 record {
    int id;
    int a = C4 + 964;
    int b = C13 * 6;
    int c = (964 + C4) % 10000;
    boolean flag = true;
    string label = "rec964";
};
type R965 record {
    int id;
    int a = C5 + 965;
    int b = C16 * 7;
    int c = (965 + C5) % 10000;
    boolean flag = false;
    string label = "rec965";
};
type R966 record {
    int id;
    int a = C6 + 966;
    int b = C19 * 1;
    int c = (966 + C6) % 10000;
    boolean flag = true;
    string label = "rec966";
};
type R967 record {
    int id;
    int a = C7 + 967;
    int b = C22 * 2;
    int c = (967 + C7) % 10000;
    boolean flag = false;
    string label = "rec967";
};
type R968 record {
    int id;
    int a = C8 + 968;
    int b = C25 * 3;
    int c = (968 + C8) % 10000;
    boolean flag = true;
    string label = "rec968";
};
type R969 record {
    int id;
    int a = C9 + 969;
    int b = C28 * 4;
    int c = (969 + C9) % 10000;
    boolean flag = false;
    string label = "rec969";
};
type R970 record {
    int id;
    int a = C10 + 970;
    int b = C31 * 5;
    int c = (970 + C10) % 10000;
    boolean flag = true;
    string label = "rec970";
};
type R971 record {
    int id;
    int a = C11 + 971;
    int b = C34 * 6;
    int c = (971 + C11) % 10000;
    boolean flag = false;
    string label = "rec971";
};
type R972 record {
    int id;
    int a = C12 + 972;
    int b = C37 * 7;
    int c = (972 + C12) % 10000;
    boolean flag = true;
    string label = "rec972";
};
type R973 record {
    int id;
    int a = C13 + 973;
    int b = C40 * 1;
    int c = (973 + C13) % 10000;
    boolean flag = false;
    string label = "rec973";
};
type R974 record {
    int id;
    int a = C14 + 974;
    int b = C43 * 2;
    int c = (974 + C14) % 10000;
    boolean flag = true;
    string label = "rec974";
};
type R975 record {
    int id;
    int a = C15 + 975;
    int b = C46 * 3;
    int c = (975 + C15) % 10000;
    boolean flag = false;
    string label = "rec975";
};
type R976 record {
    int id;
    int a = C16 + 976;
    int b = C49 * 4;
    int c = (976 + C16) % 10000;
    boolean flag = true;
    string label = "rec976";
};
type R977 record {
    int id;
    int a = C17 + 977;
    int b = C52 * 5;
    int c = (977 + C17) % 10000;
    boolean flag = false;
    string label = "rec977";
};
type R978 record {
    int id;
    int a = C18 + 978;
    int b = C55 * 6;
    int c = (978 + C18) % 10000;
    boolean flag = true;
    string label = "rec978";
};
type R979 record {
    int id;
    int a = C19 + 979;
    int b = C58 * 7;
    int c = (979 + C19) % 10000;
    boolean flag = false;
    string label = "rec979";
};
type R980 record {
    int id;
    int a = C20 + 980;
    int b = C61 * 1;
    int c = (980 + C20) % 10000;
    boolean flag = true;
    string label = "rec980";
};
type R981 record {
    int id;
    int a = C21 + 981;
    int b = C0 * 2;
    int c = (981 + C21) % 10000;
    boolean flag = false;
    string label = "rec981";
};
type R982 record {
    int id;
    int a = C22 + 982;
    int b = C3 * 3;
    int c = (982 + C22) % 10000;
    boolean flag = true;
    string label = "rec982";
};
type R983 record {
    int id;
    int a = C23 + 983;
    int b = C6 * 4;
    int c = (983 + C23) % 10000;
    boolean flag = false;
    string label = "rec983";
};
type R984 record {
    int id;
    int a = C24 + 984;
    int b = C9 * 5;
    int c = (984 + C24) % 10000;
    boolean flag = true;
    string label = "rec984";
};
type R985 record {
    int id;
    int a = C25 + 985;
    int b = C12 * 6;
    int c = (985 + C25) % 10000;
    boolean flag = false;
    string label = "rec985";
};
type R986 record {
    int id;
    int a = C26 + 986;
    int b = C15 * 7;
    int c = (986 + C26) % 10000;
    boolean flag = true;
    string label = "rec986";
};
type R987 record {
    int id;
    int a = C27 + 987;
    int b = C18 * 1;
    int c = (987 + C27) % 10000;
    boolean flag = false;
    string label = "rec987";
};
type R988 record {
    int id;
    int a = C28 + 988;
    int b = C21 * 2;
    int c = (988 + C28) % 10000;
    boolean flag = true;
    string label = "rec988";
};
type R989 record {
    int id;
    int a = C29 + 989;
    int b = C24 * 3;
    int c = (989 + C29) % 10000;
    boolean flag = false;
    string label = "rec989";
};
type R990 record {
    int id;
    int a = C30 + 990;
    int b = C27 * 4;
    int c = (990 + C30) % 10000;
    boolean flag = true;
    string label = "rec990";
};
type R991 record {
    int id;
    int a = C31 + 991;
    int b = C30 * 5;
    int c = (991 + C31) % 10000;
    boolean flag = false;
    string label = "rec991";
};
type R992 record {
    int id;
    int a = C32 + 992;
    int b = C33 * 6;
    int c = (992 + C32) % 10000;
    boolean flag = true;
    string label = "rec992";
};
type R993 record {
    int id;
    int a = C33 + 993;
    int b = C36 * 7;
    int c = (993 + C33) % 10000;
    boolean flag = false;
    string label = "rec993";
};
type R994 record {
    int id;
    int a = C34 + 994;
    int b = C39 * 1;
    int c = (994 + C34) % 10000;
    boolean flag = true;
    string label = "rec994";
};
type R995 record {
    int id;
    int a = C35 + 995;
    int b = C42 * 2;
    int c = (995 + C35) % 10000;
    boolean flag = false;
    string label = "rec995";
};
type R996 record {
    int id;
    int a = C36 + 996;
    int b = C45 * 3;
    int c = (996 + C36) % 10000;
    boolean flag = true;
    string label = "rec996";
};
type R997 record {
    int id;
    int a = C37 + 997;
    int b = C48 * 4;
    int c = (997 + C37) % 10000;
    boolean flag = false;
    string label = "rec997";
};
type R998 record {
    int id;
    int a = C38 + 998;
    int b = C51 * 5;
    int c = (998 + C38) % 10000;
    boolean flag = true;
    string label = "rec998";
};
type R999 record {
    int id;
    int a = C39 + 999;
    int b = C54 * 6;
    int c = (999 + C39) % 10000;
    boolean flag = false;
    string label = "rec999";
};
type R1000 record {
    int id;
    int a = C40 + 1000;
    int b = C57 * 7;
    int c = (1000 + C40) % 10000;
    boolean flag = true;
    string label = "rec1000";
};
type R1001 record {
    int id;
    int a = C41 + 1001;
    int b = C60 * 1;
    int c = (1001 + C41) % 10000;
    boolean flag = false;
    string label = "rec1001";
};
type R1002 record {
    int id;
    int a = C42 + 1002;
    int b = C63 * 2;
    int c = (1002 + C42) % 10000;
    boolean flag = true;
    string label = "rec1002";
};
type R1003 record {
    int id;
    int a = C43 + 1003;
    int b = C2 * 3;
    int c = (1003 + C43) % 10000;
    boolean flag = false;
    string label = "rec1003";
};
type R1004 record {
    int id;
    int a = C44 + 1004;
    int b = C5 * 4;
    int c = (1004 + C44) % 10000;
    boolean flag = true;
    string label = "rec1004";
};
type R1005 record {
    int id;
    int a = C45 + 1005;
    int b = C8 * 5;
    int c = (1005 + C45) % 10000;
    boolean flag = false;
    string label = "rec1005";
};
type R1006 record {
    int id;
    int a = C46 + 1006;
    int b = C11 * 6;
    int c = (1006 + C46) % 10000;
    boolean flag = true;
    string label = "rec1006";
};
type R1007 record {
    int id;
    int a = C47 + 1007;
    int b = C14 * 7;
    int c = (1007 + C47) % 10000;
    boolean flag = false;
    string label = "rec1007";
};
type R1008 record {
    int id;
    int a = C48 + 1008;
    int b = C17 * 1;
    int c = (1008 + C48) % 10000;
    boolean flag = true;
    string label = "rec1008";
};
type R1009 record {
    int id;
    int a = C49 + 1009;
    int b = C20 * 2;
    int c = (1009 + C49) % 10000;
    boolean flag = false;
    string label = "rec1009";
};
type R1010 record {
    int id;
    int a = C50 + 1010;
    int b = C23 * 3;
    int c = (1010 + C50) % 10000;
    boolean flag = true;
    string label = "rec1010";
};
type R1011 record {
    int id;
    int a = C51 + 1011;
    int b = C26 * 4;
    int c = (1011 + C51) % 10000;
    boolean flag = false;
    string label = "rec1011";
};
type R1012 record {
    int id;
    int a = C52 + 1012;
    int b = C29 * 5;
    int c = (1012 + C52) % 10000;
    boolean flag = true;
    string label = "rec1012";
};
type R1013 record {
    int id;
    int a = C53 + 1013;
    int b = C32 * 6;
    int c = (1013 + C53) % 10000;
    boolean flag = false;
    string label = "rec1013";
};
type R1014 record {
    int id;
    int a = C54 + 1014;
    int b = C35 * 7;
    int c = (1014 + C54) % 10000;
    boolean flag = true;
    string label = "rec1014";
};
type R1015 record {
    int id;
    int a = C55 + 1015;
    int b = C38 * 1;
    int c = (1015 + C55) % 10000;
    boolean flag = false;
    string label = "rec1015";
};
type R1016 record {
    int id;
    int a = C56 + 1016;
    int b = C41 * 2;
    int c = (1016 + C56) % 10000;
    boolean flag = true;
    string label = "rec1016";
};
type R1017 record {
    int id;
    int a = C57 + 1017;
    int b = C44 * 3;
    int c = (1017 + C57) % 10000;
    boolean flag = false;
    string label = "rec1017";
};
type R1018 record {
    int id;
    int a = C58 + 1018;
    int b = C47 * 4;
    int c = (1018 + C58) % 10000;
    boolean flag = true;
    string label = "rec1018";
};
type R1019 record {
    int id;
    int a = C59 + 1019;
    int b = C50 * 5;
    int c = (1019 + C59) % 10000;
    boolean flag = false;
    string label = "rec1019";
};
type R1020 record {
    int id;
    int a = C60 + 1020;
    int b = C53 * 6;
    int c = (1020 + C60) % 10000;
    boolean flag = true;
    string label = "rec1020";
};
type R1021 record {
    int id;
    int a = C61 + 1021;
    int b = C56 * 7;
    int c = (1021 + C61) % 10000;
    boolean flag = false;
    string label = "rec1021";
};
type R1022 record {
    int id;
    int a = C62 + 1022;
    int b = C59 * 1;
    int c = (1022 + C62) % 10000;
    boolean flag = true;
    string label = "rec1022";
};
type R1023 record {
    int id;
    int a = C63 + 1023;
    int b = C62 * 2;
    int c = (1023 + C63) % 10000;
    boolean flag = false;
    string label = "rec1023";
};
type R1024 record {
    int id;
    int a = C0 + 1024;
    int b = C1 * 3;
    int c = (1024 + C0) % 10000;
    boolean flag = true;
    string label = "rec1024";
};
type R1025 record {
    int id;
    int a = C1 + 1025;
    int b = C4 * 4;
    int c = (1025 + C1) % 10000;
    boolean flag = false;
    string label = "rec1025";
};
type R1026 record {
    int id;
    int a = C2 + 1026;
    int b = C7 * 5;
    int c = (1026 + C2) % 10000;
    boolean flag = true;
    string label = "rec1026";
};
type R1027 record {
    int id;
    int a = C3 + 1027;
    int b = C10 * 6;
    int c = (1027 + C3) % 10000;
    boolean flag = false;
    string label = "rec1027";
};
type R1028 record {
    int id;
    int a = C4 + 1028;
    int b = C13 * 7;
    int c = (1028 + C4) % 10000;
    boolean flag = true;
    string label = "rec1028";
};
type R1029 record {
    int id;
    int a = C5 + 1029;
    int b = C16 * 1;
    int c = (1029 + C5) % 10000;
    boolean flag = false;
    string label = "rec1029";
};
type R1030 record {
    int id;
    int a = C6 + 1030;
    int b = C19 * 2;
    int c = (1030 + C6) % 10000;
    boolean flag = true;
    string label = "rec1030";
};
type R1031 record {
    int id;
    int a = C7 + 1031;
    int b = C22 * 3;
    int c = (1031 + C7) % 10000;
    boolean flag = false;
    string label = "rec1031";
};
type R1032 record {
    int id;
    int a = C8 + 1032;
    int b = C25 * 4;
    int c = (1032 + C8) % 10000;
    boolean flag = true;
    string label = "rec1032";
};
type R1033 record {
    int id;
    int a = C9 + 1033;
    int b = C28 * 5;
    int c = (1033 + C9) % 10000;
    boolean flag = false;
    string label = "rec1033";
};
type R1034 record {
    int id;
    int a = C10 + 1034;
    int b = C31 * 6;
    int c = (1034 + C10) % 10000;
    boolean flag = true;
    string label = "rec1034";
};
type R1035 record {
    int id;
    int a = C11 + 1035;
    int b = C34 * 7;
    int c = (1035 + C11) % 10000;
    boolean flag = false;
    string label = "rec1035";
};
type R1036 record {
    int id;
    int a = C12 + 1036;
    int b = C37 * 1;
    int c = (1036 + C12) % 10000;
    boolean flag = true;
    string label = "rec1036";
};
type R1037 record {
    int id;
    int a = C13 + 1037;
    int b = C40 * 2;
    int c = (1037 + C13) % 10000;
    boolean flag = false;
    string label = "rec1037";
};
type R1038 record {
    int id;
    int a = C14 + 1038;
    int b = C43 * 3;
    int c = (1038 + C14) % 10000;
    boolean flag = true;
    string label = "rec1038";
};
type R1039 record {
    int id;
    int a = C15 + 1039;
    int b = C46 * 4;
    int c = (1039 + C15) % 10000;
    boolean flag = false;
    string label = "rec1039";
};
type R1040 record {
    int id;
    int a = C16 + 1040;
    int b = C49 * 5;
    int c = (1040 + C16) % 10000;
    boolean flag = true;
    string label = "rec1040";
};
type R1041 record {
    int id;
    int a = C17 + 1041;
    int b = C52 * 6;
    int c = (1041 + C17) % 10000;
    boolean flag = false;
    string label = "rec1041";
};
type R1042 record {
    int id;
    int a = C18 + 1042;
    int b = C55 * 7;
    int c = (1042 + C18) % 10000;
    boolean flag = true;
    string label = "rec1042";
};
type R1043 record {
    int id;
    int a = C19 + 1043;
    int b = C58 * 1;
    int c = (1043 + C19) % 10000;
    boolean flag = false;
    string label = "rec1043";
};
type R1044 record {
    int id;
    int a = C20 + 1044;
    int b = C61 * 2;
    int c = (1044 + C20) % 10000;
    boolean flag = true;
    string label = "rec1044";
};
type R1045 record {
    int id;
    int a = C21 + 1045;
    int b = C0 * 3;
    int c = (1045 + C21) % 10000;
    boolean flag = false;
    string label = "rec1045";
};
type R1046 record {
    int id;
    int a = C22 + 1046;
    int b = C3 * 4;
    int c = (1046 + C22) % 10000;
    boolean flag = true;
    string label = "rec1046";
};
type R1047 record {
    int id;
    int a = C23 + 1047;
    int b = C6 * 5;
    int c = (1047 + C23) % 10000;
    boolean flag = false;
    string label = "rec1047";
};
type R1048 record {
    int id;
    int a = C24 + 1048;
    int b = C9 * 6;
    int c = (1048 + C24) % 10000;
    boolean flag = true;
    string label = "rec1048";
};
type R1049 record {
    int id;
    int a = C25 + 1049;
    int b = C12 * 7;
    int c = (1049 + C25) % 10000;
    boolean flag = false;
    string label = "rec1049";
};
type R1050 record {
    int id;
    int a = C26 + 1050;
    int b = C15 * 1;
    int c = (1050 + C26) % 10000;
    boolean flag = true;
    string label = "rec1050";
};
type R1051 record {
    int id;
    int a = C27 + 1051;
    int b = C18 * 2;
    int c = (1051 + C27) % 10000;
    boolean flag = false;
    string label = "rec1051";
};
type R1052 record {
    int id;
    int a = C28 + 1052;
    int b = C21 * 3;
    int c = (1052 + C28) % 10000;
    boolean flag = true;
    string label = "rec1052";
};
type R1053 record {
    int id;
    int a = C29 + 1053;
    int b = C24 * 4;
    int c = (1053 + C29) % 10000;
    boolean flag = false;
    string label = "rec1053";
};
type R1054 record {
    int id;
    int a = C30 + 1054;
    int b = C27 * 5;
    int c = (1054 + C30) % 10000;
    boolean flag = true;
    string label = "rec1054";
};
type R1055 record {
    int id;
    int a = C31 + 1055;
    int b = C30 * 6;
    int c = (1055 + C31) % 10000;
    boolean flag = false;
    string label = "rec1055";
};
type R1056 record {
    int id;
    int a = C32 + 1056;
    int b = C33 * 7;
    int c = (1056 + C32) % 10000;
    boolean flag = true;
    string label = "rec1056";
};
type R1057 record {
    int id;
    int a = C33 + 1057;
    int b = C36 * 1;
    int c = (1057 + C33) % 10000;
    boolean flag = false;
    string label = "rec1057";
};
type R1058 record {
    int id;
    int a = C34 + 1058;
    int b = C39 * 2;
    int c = (1058 + C34) % 10000;
    boolean flag = true;
    string label = "rec1058";
};
type R1059 record {
    int id;
    int a = C35 + 1059;
    int b = C42 * 3;
    int c = (1059 + C35) % 10000;
    boolean flag = false;
    string label = "rec1059";
};
type R1060 record {
    int id;
    int a = C36 + 1060;
    int b = C45 * 4;
    int c = (1060 + C36) % 10000;
    boolean flag = true;
    string label = "rec1060";
};
type R1061 record {
    int id;
    int a = C37 + 1061;
    int b = C48 * 5;
    int c = (1061 + C37) % 10000;
    boolean flag = false;
    string label = "rec1061";
};
type R1062 record {
    int id;
    int a = C38 + 1062;
    int b = C51 * 6;
    int c = (1062 + C38) % 10000;
    boolean flag = true;
    string label = "rec1062";
};
type R1063 record {
    int id;
    int a = C39 + 1063;
    int b = C54 * 7;
    int c = (1063 + C39) % 10000;
    boolean flag = false;
    string label = "rec1063";
};
type R1064 record {
    int id;
    int a = C40 + 1064;
    int b = C57 * 1;
    int c = (1064 + C40) % 10000;
    boolean flag = true;
    string label = "rec1064";
};
type R1065 record {
    int id;
    int a = C41 + 1065;
    int b = C60 * 2;
    int c = (1065 + C41) % 10000;
    boolean flag = false;
    string label = "rec1065";
};
type R1066 record {
    int id;
    int a = C42 + 1066;
    int b = C63 * 3;
    int c = (1066 + C42) % 10000;
    boolean flag = true;
    string label = "rec1066";
};
type R1067 record {
    int id;
    int a = C43 + 1067;
    int b = C2 * 4;
    int c = (1067 + C43) % 10000;
    boolean flag = false;
    string label = "rec1067";
};
type R1068 record {
    int id;
    int a = C44 + 1068;
    int b = C5 * 5;
    int c = (1068 + C44) % 10000;
    boolean flag = true;
    string label = "rec1068";
};
type R1069 record {
    int id;
    int a = C45 + 1069;
    int b = C8 * 6;
    int c = (1069 + C45) % 10000;
    boolean flag = false;
    string label = "rec1069";
};
type R1070 record {
    int id;
    int a = C46 + 1070;
    int b = C11 * 7;
    int c = (1070 + C46) % 10000;
    boolean flag = true;
    string label = "rec1070";
};
type R1071 record {
    int id;
    int a = C47 + 1071;
    int b = C14 * 1;
    int c = (1071 + C47) % 10000;
    boolean flag = false;
    string label = "rec1071";
};
type R1072 record {
    int id;
    int a = C48 + 1072;
    int b = C17 * 2;
    int c = (1072 + C48) % 10000;
    boolean flag = true;
    string label = "rec1072";
};
type R1073 record {
    int id;
    int a = C49 + 1073;
    int b = C20 * 3;
    int c = (1073 + C49) % 10000;
    boolean flag = false;
    string label = "rec1073";
};
type R1074 record {
    int id;
    int a = C50 + 1074;
    int b = C23 * 4;
    int c = (1074 + C50) % 10000;
    boolean flag = true;
    string label = "rec1074";
};
type R1075 record {
    int id;
    int a = C51 + 1075;
    int b = C26 * 5;
    int c = (1075 + C51) % 10000;
    boolean flag = false;
    string label = "rec1075";
};
type R1076 record {
    int id;
    int a = C52 + 1076;
    int b = C29 * 6;
    int c = (1076 + C52) % 10000;
    boolean flag = true;
    string label = "rec1076";
};
type R1077 record {
    int id;
    int a = C53 + 1077;
    int b = C32 * 7;
    int c = (1077 + C53) % 10000;
    boolean flag = false;
    string label = "rec1077";
};
type R1078 record {
    int id;
    int a = C54 + 1078;
    int b = C35 * 1;
    int c = (1078 + C54) % 10000;
    boolean flag = true;
    string label = "rec1078";
};
type R1079 record {
    int id;
    int a = C55 + 1079;
    int b = C38 * 2;
    int c = (1079 + C55) % 10000;
    boolean flag = false;
    string label = "rec1079";
};
type R1080 record {
    int id;
    int a = C56 + 1080;
    int b = C41 * 3;
    int c = (1080 + C56) % 10000;
    boolean flag = true;
    string label = "rec1080";
};
type R1081 record {
    int id;
    int a = C57 + 1081;
    int b = C44 * 4;
    int c = (1081 + C57) % 10000;
    boolean flag = false;
    string label = "rec1081";
};
type R1082 record {
    int id;
    int a = C58 + 1082;
    int b = C47 * 5;
    int c = (1082 + C58) % 10000;
    boolean flag = true;
    string label = "rec1082";
};
type R1083 record {
    int id;
    int a = C59 + 1083;
    int b = C50 * 6;
    int c = (1083 + C59) % 10000;
    boolean flag = false;
    string label = "rec1083";
};
type R1084 record {
    int id;
    int a = C60 + 1084;
    int b = C53 * 7;
    int c = (1084 + C60) % 10000;
    boolean flag = true;
    string label = "rec1084";
};
type R1085 record {
    int id;
    int a = C61 + 1085;
    int b = C56 * 1;
    int c = (1085 + C61) % 10000;
    boolean flag = false;
    string label = "rec1085";
};
type R1086 record {
    int id;
    int a = C62 + 1086;
    int b = C59 * 2;
    int c = (1086 + C62) % 10000;
    boolean flag = true;
    string label = "rec1086";
};
type R1087 record {
    int id;
    int a = C63 + 1087;
    int b = C62 * 3;
    int c = (1087 + C63) % 10000;
    boolean flag = false;
    string label = "rec1087";
};
type R1088 record {
    int id;
    int a = C0 + 1088;
    int b = C1 * 4;
    int c = (1088 + C0) % 10000;
    boolean flag = true;
    string label = "rec1088";
};
type R1089 record {
    int id;
    int a = C1 + 1089;
    int b = C4 * 5;
    int c = (1089 + C1) % 10000;
    boolean flag = false;
    string label = "rec1089";
};
type R1090 record {
    int id;
    int a = C2 + 1090;
    int b = C7 * 6;
    int c = (1090 + C2) % 10000;
    boolean flag = true;
    string label = "rec1090";
};
type R1091 record {
    int id;
    int a = C3 + 1091;
    int b = C10 * 7;
    int c = (1091 + C3) % 10000;
    boolean flag = false;
    string label = "rec1091";
};
type R1092 record {
    int id;
    int a = C4 + 1092;
    int b = C13 * 1;
    int c = (1092 + C4) % 10000;
    boolean flag = true;
    string label = "rec1092";
};
type R1093 record {
    int id;
    int a = C5 + 1093;
    int b = C16 * 2;
    int c = (1093 + C5) % 10000;
    boolean flag = false;
    string label = "rec1093";
};
type R1094 record {
    int id;
    int a = C6 + 1094;
    int b = C19 * 3;
    int c = (1094 + C6) % 10000;
    boolean flag = true;
    string label = "rec1094";
};
type R1095 record {
    int id;
    int a = C7 + 1095;
    int b = C22 * 4;
    int c = (1095 + C7) % 10000;
    boolean flag = false;
    string label = "rec1095";
};
type R1096 record {
    int id;
    int a = C8 + 1096;
    int b = C25 * 5;
    int c = (1096 + C8) % 10000;
    boolean flag = true;
    string label = "rec1096";
};
type R1097 record {
    int id;
    int a = C9 + 1097;
    int b = C28 * 6;
    int c = (1097 + C9) % 10000;
    boolean flag = false;
    string label = "rec1097";
};
type R1098 record {
    int id;
    int a = C10 + 1098;
    int b = C31 * 7;
    int c = (1098 + C10) % 10000;
    boolean flag = true;
    string label = "rec1098";
};
type R1099 record {
    int id;
    int a = C11 + 1099;
    int b = C34 * 1;
    int c = (1099 + C11) % 10000;
    boolean flag = false;
    string label = "rec1099";
};
type R1100 record {
    int id;
    int a = C12 + 1100;
    int b = C37 * 2;
    int c = (1100 + C12) % 10000;
    boolean flag = true;
    string label = "rec1100";
};
type R1101 record {
    int id;
    int a = C13 + 1101;
    int b = C40 * 3;
    int c = (1101 + C13) % 10000;
    boolean flag = false;
    string label = "rec1101";
};
type R1102 record {
    int id;
    int a = C14 + 1102;
    int b = C43 * 4;
    int c = (1102 + C14) % 10000;
    boolean flag = true;
    string label = "rec1102";
};
type R1103 record {
    int id;
    int a = C15 + 1103;
    int b = C46 * 5;
    int c = (1103 + C15) % 10000;
    boolean flag = false;
    string label = "rec1103";
};
type R1104 record {
    int id;
    int a = C16 + 1104;
    int b = C49 * 6;
    int c = (1104 + C16) % 10000;
    boolean flag = true;
    string label = "rec1104";
};
type R1105 record {
    int id;
    int a = C17 + 1105;
    int b = C52 * 7;
    int c = (1105 + C17) % 10000;
    boolean flag = false;
    string label = "rec1105";
};
type R1106 record {
    int id;
    int a = C18 + 1106;
    int b = C55 * 1;
    int c = (1106 + C18) % 10000;
    boolean flag = true;
    string label = "rec1106";
};
type R1107 record {
    int id;
    int a = C19 + 1107;
    int b = C58 * 2;
    int c = (1107 + C19) % 10000;
    boolean flag = false;
    string label = "rec1107";
};
type R1108 record {
    int id;
    int a = C20 + 1108;
    int b = C61 * 3;
    int c = (1108 + C20) % 10000;
    boolean flag = true;
    string label = "rec1108";
};
type R1109 record {
    int id;
    int a = C21 + 1109;
    int b = C0 * 4;
    int c = (1109 + C21) % 10000;
    boolean flag = false;
    string label = "rec1109";
};
type R1110 record {
    int id;
    int a = C22 + 1110;
    int b = C3 * 5;
    int c = (1110 + C22) % 10000;
    boolean flag = true;
    string label = "rec1110";
};
type R1111 record {
    int id;
    int a = C23 + 1111;
    int b = C6 * 6;
    int c = (1111 + C23) % 10000;
    boolean flag = false;
    string label = "rec1111";
};
type R1112 record {
    int id;
    int a = C24 + 1112;
    int b = C9 * 7;
    int c = (1112 + C24) % 10000;
    boolean flag = true;
    string label = "rec1112";
};
type R1113 record {
    int id;
    int a = C25 + 1113;
    int b = C12 * 1;
    int c = (1113 + C25) % 10000;
    boolean flag = false;
    string label = "rec1113";
};
type R1114 record {
    int id;
    int a = C26 + 1114;
    int b = C15 * 2;
    int c = (1114 + C26) % 10000;
    boolean flag = true;
    string label = "rec1114";
};
type R1115 record {
    int id;
    int a = C27 + 1115;
    int b = C18 * 3;
    int c = (1115 + C27) % 10000;
    boolean flag = false;
    string label = "rec1115";
};
type R1116 record {
    int id;
    int a = C28 + 1116;
    int b = C21 * 4;
    int c = (1116 + C28) % 10000;
    boolean flag = true;
    string label = "rec1116";
};
type R1117 record {
    int id;
    int a = C29 + 1117;
    int b = C24 * 5;
    int c = (1117 + C29) % 10000;
    boolean flag = false;
    string label = "rec1117";
};
type R1118 record {
    int id;
    int a = C30 + 1118;
    int b = C27 * 6;
    int c = (1118 + C30) % 10000;
    boolean flag = true;
    string label = "rec1118";
};
type R1119 record {
    int id;
    int a = C31 + 1119;
    int b = C30 * 7;
    int c = (1119 + C31) % 10000;
    boolean flag = false;
    string label = "rec1119";
};
type R1120 record {
    int id;
    int a = C32 + 1120;
    int b = C33 * 1;
    int c = (1120 + C32) % 10000;
    boolean flag = true;
    string label = "rec1120";
};
type R1121 record {
    int id;
    int a = C33 + 1121;
    int b = C36 * 2;
    int c = (1121 + C33) % 10000;
    boolean flag = false;
    string label = "rec1121";
};
type R1122 record {
    int id;
    int a = C34 + 1122;
    int b = C39 * 3;
    int c = (1122 + C34) % 10000;
    boolean flag = true;
    string label = "rec1122";
};
type R1123 record {
    int id;
    int a = C35 + 1123;
    int b = C42 * 4;
    int c = (1123 + C35) % 10000;
    boolean flag = false;
    string label = "rec1123";
};
type R1124 record {
    int id;
    int a = C36 + 1124;
    int b = C45 * 5;
    int c = (1124 + C36) % 10000;
    boolean flag = true;
    string label = "rec1124";
};
type R1125 record {
    int id;
    int a = C37 + 1125;
    int b = C48 * 6;
    int c = (1125 + C37) % 10000;
    boolean flag = false;
    string label = "rec1125";
};
type R1126 record {
    int id;
    int a = C38 + 1126;
    int b = C51 * 7;
    int c = (1126 + C38) % 10000;
    boolean flag = true;
    string label = "rec1126";
};
type R1127 record {
    int id;
    int a = C39 + 1127;
    int b = C54 * 1;
    int c = (1127 + C39) % 10000;
    boolean flag = false;
    string label = "rec1127";
};
type R1128 record {
    int id;
    int a = C40 + 1128;
    int b = C57 * 2;
    int c = (1128 + C40) % 10000;
    boolean flag = true;
    string label = "rec1128";
};
type R1129 record {
    int id;
    int a = C41 + 1129;
    int b = C60 * 3;
    int c = (1129 + C41) % 10000;
    boolean flag = false;
    string label = "rec1129";
};
type R1130 record {
    int id;
    int a = C42 + 1130;
    int b = C63 * 4;
    int c = (1130 + C42) % 10000;
    boolean flag = true;
    string label = "rec1130";
};
type R1131 record {
    int id;
    int a = C43 + 1131;
    int b = C2 * 5;
    int c = (1131 + C43) % 10000;
    boolean flag = false;
    string label = "rec1131";
};
type R1132 record {
    int id;
    int a = C44 + 1132;
    int b = C5 * 6;
    int c = (1132 + C44) % 10000;
    boolean flag = true;
    string label = "rec1132";
};
type R1133 record {
    int id;
    int a = C45 + 1133;
    int b = C8 * 7;
    int c = (1133 + C45) % 10000;
    boolean flag = false;
    string label = "rec1133";
};
type R1134 record {
    int id;
    int a = C46 + 1134;
    int b = C11 * 1;
    int c = (1134 + C46) % 10000;
    boolean flag = true;
    string label = "rec1134";
};
type R1135 record {
    int id;
    int a = C47 + 1135;
    int b = C14 * 2;
    int c = (1135 + C47) % 10000;
    boolean flag = false;
    string label = "rec1135";
};
type R1136 record {
    int id;
    int a = C48 + 1136;
    int b = C17 * 3;
    int c = (1136 + C48) % 10000;
    boolean flag = true;
    string label = "rec1136";
};
type R1137 record {
    int id;
    int a = C49 + 1137;
    int b = C20 * 4;
    int c = (1137 + C49) % 10000;
    boolean flag = false;
    string label = "rec1137";
};
type R1138 record {
    int id;
    int a = C50 + 1138;
    int b = C23 * 5;
    int c = (1138 + C50) % 10000;
    boolean flag = true;
    string label = "rec1138";
};
type R1139 record {
    int id;
    int a = C51 + 1139;
    int b = C26 * 6;
    int c = (1139 + C51) % 10000;
    boolean flag = false;
    string label = "rec1139";
};
type R1140 record {
    int id;
    int a = C52 + 1140;
    int b = C29 * 7;
    int c = (1140 + C52) % 10000;
    boolean flag = true;
    string label = "rec1140";
};
type R1141 record {
    int id;
    int a = C53 + 1141;
    int b = C32 * 1;
    int c = (1141 + C53) % 10000;
    boolean flag = false;
    string label = "rec1141";
};
type R1142 record {
    int id;
    int a = C54 + 1142;
    int b = C35 * 2;
    int c = (1142 + C54) % 10000;
    boolean flag = true;
    string label = "rec1142";
};
type R1143 record {
    int id;
    int a = C55 + 1143;
    int b = C38 * 3;
    int c = (1143 + C55) % 10000;
    boolean flag = false;
    string label = "rec1143";
};
type R1144 record {
    int id;
    int a = C56 + 1144;
    int b = C41 * 4;
    int c = (1144 + C56) % 10000;
    boolean flag = true;
    string label = "rec1144";
};
type R1145 record {
    int id;
    int a = C57 + 1145;
    int b = C44 * 5;
    int c = (1145 + C57) % 10000;
    boolean flag = false;
    string label = "rec1145";
};
type R1146 record {
    int id;
    int a = C58 + 1146;
    int b = C47 * 6;
    int c = (1146 + C58) % 10000;
    boolean flag = true;
    string label = "rec1146";
};
type R1147 record {
    int id;
    int a = C59 + 1147;
    int b = C50 * 7;
    int c = (1147 + C59) % 10000;
    boolean flag = false;
    string label = "rec1147";
};
type R1148 record {
    int id;
    int a = C60 + 1148;
    int b = C53 * 1;
    int c = (1148 + C60) % 10000;
    boolean flag = true;
    string label = "rec1148";
};
type R1149 record {
    int id;
    int a = C61 + 1149;
    int b = C56 * 2;
    int c = (1149 + C61) % 10000;
    boolean flag = false;
    string label = "rec1149";
};
type R1150 record {
    int id;
    int a = C62 + 1150;
    int b = C59 * 3;
    int c = (1150 + C62) % 10000;
    boolean flag = true;
    string label = "rec1150";
};
type R1151 record {
    int id;
    int a = C63 + 1151;
    int b = C62 * 4;
    int c = (1151 + C63) % 10000;
    boolean flag = false;
    string label = "rec1151";
};
type R1152 record {
    int id;
    int a = C0 + 1152;
    int b = C1 * 5;
    int c = (1152 + C0) % 10000;
    boolean flag = true;
    string label = "rec1152";
};
type R1153 record {
    int id;
    int a = C1 + 1153;
    int b = C4 * 6;
    int c = (1153 + C1) % 10000;
    boolean flag = false;
    string label = "rec1153";
};
type R1154 record {
    int id;
    int a = C2 + 1154;
    int b = C7 * 7;
    int c = (1154 + C2) % 10000;
    boolean flag = true;
    string label = "rec1154";
};
type R1155 record {
    int id;
    int a = C3 + 1155;
    int b = C10 * 1;
    int c = (1155 + C3) % 10000;
    boolean flag = false;
    string label = "rec1155";
};
type R1156 record {
    int id;
    int a = C4 + 1156;
    int b = C13 * 2;
    int c = (1156 + C4) % 10000;
    boolean flag = true;
    string label = "rec1156";
};
type R1157 record {
    int id;
    int a = C5 + 1157;
    int b = C16 * 3;
    int c = (1157 + C5) % 10000;
    boolean flag = false;
    string label = "rec1157";
};
type R1158 record {
    int id;
    int a = C6 + 1158;
    int b = C19 * 4;
    int c = (1158 + C6) % 10000;
    boolean flag = true;
    string label = "rec1158";
};
type R1159 record {
    int id;
    int a = C7 + 1159;
    int b = C22 * 5;
    int c = (1159 + C7) % 10000;
    boolean flag = false;
    string label = "rec1159";
};
type R1160 record {
    int id;
    int a = C8 + 1160;
    int b = C25 * 6;
    int c = (1160 + C8) % 10000;
    boolean flag = true;
    string label = "rec1160";
};
type R1161 record {
    int id;
    int a = C9 + 1161;
    int b = C28 * 7;
    int c = (1161 + C9) % 10000;
    boolean flag = false;
    string label = "rec1161";
};
type R1162 record {
    int id;
    int a = C10 + 1162;
    int b = C31 * 1;
    int c = (1162 + C10) % 10000;
    boolean flag = true;
    string label = "rec1162";
};
type R1163 record {
    int id;
    int a = C11 + 1163;
    int b = C34 * 2;
    int c = (1163 + C11) % 10000;
    boolean flag = false;
    string label = "rec1163";
};
type R1164 record {
    int id;
    int a = C12 + 1164;
    int b = C37 * 3;
    int c = (1164 + C12) % 10000;
    boolean flag = true;
    string label = "rec1164";
};
type R1165 record {
    int id;
    int a = C13 + 1165;
    int b = C40 * 4;
    int c = (1165 + C13) % 10000;
    boolean flag = false;
    string label = "rec1165";
};
type R1166 record {
    int id;
    int a = C14 + 1166;
    int b = C43 * 5;
    int c = (1166 + C14) % 10000;
    boolean flag = true;
    string label = "rec1166";
};
type R1167 record {
    int id;
    int a = C15 + 1167;
    int b = C46 * 6;
    int c = (1167 + C15) % 10000;
    boolean flag = false;
    string label = "rec1167";
};
type R1168 record {
    int id;
    int a = C16 + 1168;
    int b = C49 * 7;
    int c = (1168 + C16) % 10000;
    boolean flag = true;
    string label = "rec1168";
};
type R1169 record {
    int id;
    int a = C17 + 1169;
    int b = C52 * 1;
    int c = (1169 + C17) % 10000;
    boolean flag = false;
    string label = "rec1169";
};
type R1170 record {
    int id;
    int a = C18 + 1170;
    int b = C55 * 2;
    int c = (1170 + C18) % 10000;
    boolean flag = true;
    string label = "rec1170";
};
type R1171 record {
    int id;
    int a = C19 + 1171;
    int b = C58 * 3;
    int c = (1171 + C19) % 10000;
    boolean flag = false;
    string label = "rec1171";
};
type R1172 record {
    int id;
    int a = C20 + 1172;
    int b = C61 * 4;
    int c = (1172 + C20) % 10000;
    boolean flag = true;
    string label = "rec1172";
};
type R1173 record {
    int id;
    int a = C21 + 1173;
    int b = C0 * 5;
    int c = (1173 + C21) % 10000;
    boolean flag = false;
    string label = "rec1173";
};
type R1174 record {
    int id;
    int a = C22 + 1174;
    int b = C3 * 6;
    int c = (1174 + C22) % 10000;
    boolean flag = true;
    string label = "rec1174";
};
type R1175 record {
    int id;
    int a = C23 + 1175;
    int b = C6 * 7;
    int c = (1175 + C23) % 10000;
    boolean flag = false;
    string label = "rec1175";
};
type R1176 record {
    int id;
    int a = C24 + 1176;
    int b = C9 * 1;
    int c = (1176 + C24) % 10000;
    boolean flag = true;
    string label = "rec1176";
};
type R1177 record {
    int id;
    int a = C25 + 1177;
    int b = C12 * 2;
    int c = (1177 + C25) % 10000;
    boolean flag = false;
    string label = "rec1177";
};
type R1178 record {
    int id;
    int a = C26 + 1178;
    int b = C15 * 3;
    int c = (1178 + C26) % 10000;
    boolean flag = true;
    string label = "rec1178";
};
type R1179 record {
    int id;
    int a = C27 + 1179;
    int b = C18 * 4;
    int c = (1179 + C27) % 10000;
    boolean flag = false;
    string label = "rec1179";
};
type R1180 record {
    int id;
    int a = C28 + 1180;
    int b = C21 * 5;
    int c = (1180 + C28) % 10000;
    boolean flag = true;
    string label = "rec1180";
};
type R1181 record {
    int id;
    int a = C29 + 1181;
    int b = C24 * 6;
    int c = (1181 + C29) % 10000;
    boolean flag = false;
    string label = "rec1181";
};
type R1182 record {
    int id;
    int a = C30 + 1182;
    int b = C27 * 7;
    int c = (1182 + C30) % 10000;
    boolean flag = true;
    string label = "rec1182";
};
type R1183 record {
    int id;
    int a = C31 + 1183;
    int b = C30 * 1;
    int c = (1183 + C31) % 10000;
    boolean flag = false;
    string label = "rec1183";
};
type R1184 record {
    int id;
    int a = C32 + 1184;
    int b = C33 * 2;
    int c = (1184 + C32) % 10000;
    boolean flag = true;
    string label = "rec1184";
};
type R1185 record {
    int id;
    int a = C33 + 1185;
    int b = C36 * 3;
    int c = (1185 + C33) % 10000;
    boolean flag = false;
    string label = "rec1185";
};
type R1186 record {
    int id;
    int a = C34 + 1186;
    int b = C39 * 4;
    int c = (1186 + C34) % 10000;
    boolean flag = true;
    string label = "rec1186";
};
type R1187 record {
    int id;
    int a = C35 + 1187;
    int b = C42 * 5;
    int c = (1187 + C35) % 10000;
    boolean flag = false;
    string label = "rec1187";
};
type R1188 record {
    int id;
    int a = C36 + 1188;
    int b = C45 * 6;
    int c = (1188 + C36) % 10000;
    boolean flag = true;
    string label = "rec1188";
};
type R1189 record {
    int id;
    int a = C37 + 1189;
    int b = C48 * 7;
    int c = (1189 + C37) % 10000;
    boolean flag = false;
    string label = "rec1189";
};
type R1190 record {
    int id;
    int a = C38 + 1190;
    int b = C51 * 1;
    int c = (1190 + C38) % 10000;
    boolean flag = true;
    string label = "rec1190";
};
type R1191 record {
    int id;
    int a = C39 + 1191;
    int b = C54 * 2;
    int c = (1191 + C39) % 10000;
    boolean flag = false;
    string label = "rec1191";
};
type R1192 record {
    int id;
    int a = C40 + 1192;
    int b = C57 * 3;
    int c = (1192 + C40) % 10000;
    boolean flag = true;
    string label = "rec1192";
};
type R1193 record {
    int id;
    int a = C41 + 1193;
    int b = C60 * 4;
    int c = (1193 + C41) % 10000;
    boolean flag = false;
    string label = "rec1193";
};
type R1194 record {
    int id;
    int a = C42 + 1194;
    int b = C63 * 5;
    int c = (1194 + C42) % 10000;
    boolean flag = true;
    string label = "rec1194";
};
type R1195 record {
    int id;
    int a = C43 + 1195;
    int b = C2 * 6;
    int c = (1195 + C43) % 10000;
    boolean flag = false;
    string label = "rec1195";
};
type R1196 record {
    int id;
    int a = C44 + 1196;
    int b = C5 * 7;
    int c = (1196 + C44) % 10000;
    boolean flag = true;
    string label = "rec1196";
};
type R1197 record {
    int id;
    int a = C45 + 1197;
    int b = C8 * 1;
    int c = (1197 + C45) % 10000;
    boolean flag = false;
    string label = "rec1197";
};
type R1198 record {
    int id;
    int a = C46 + 1198;
    int b = C11 * 2;
    int c = (1198 + C46) % 10000;
    boolean flag = true;
    string label = "rec1198";
};
type R1199 record {
    int id;
    int a = C47 + 1199;
    int b = C14 * 3;
    int c = (1199 + C47) % 10000;
    boolean flag = false;
    string label = "rec1199";
};
type R1200 record {
    int id;
    int a = C48 + 1200;
    int b = C17 * 4;
    int c = (1200 + C48) % 10000;
    boolean flag = true;
    string label = "rec1200";
};
type R1201 record {
    int id;
    int a = C49 + 1201;
    int b = C20 * 5;
    int c = (1201 + C49) % 10000;
    boolean flag = false;
    string label = "rec1201";
};
type R1202 record {
    int id;
    int a = C50 + 1202;
    int b = C23 * 6;
    int c = (1202 + C50) % 10000;
    boolean flag = true;
    string label = "rec1202";
};
type R1203 record {
    int id;
    int a = C51 + 1203;
    int b = C26 * 7;
    int c = (1203 + C51) % 10000;
    boolean flag = false;
    string label = "rec1203";
};
type R1204 record {
    int id;
    int a = C52 + 1204;
    int b = C29 * 1;
    int c = (1204 + C52) % 10000;
    boolean flag = true;
    string label = "rec1204";
};
type R1205 record {
    int id;
    int a = C53 + 1205;
    int b = C32 * 2;
    int c = (1205 + C53) % 10000;
    boolean flag = false;
    string label = "rec1205";
};
type R1206 record {
    int id;
    int a = C54 + 1206;
    int b = C35 * 3;
    int c = (1206 + C54) % 10000;
    boolean flag = true;
    string label = "rec1206";
};
type R1207 record {
    int id;
    int a = C55 + 1207;
    int b = C38 * 4;
    int c = (1207 + C55) % 10000;
    boolean flag = false;
    string label = "rec1207";
};
type R1208 record {
    int id;
    int a = C56 + 1208;
    int b = C41 * 5;
    int c = (1208 + C56) % 10000;
    boolean flag = true;
    string label = "rec1208";
};
type R1209 record {
    int id;
    int a = C57 + 1209;
    int b = C44 * 6;
    int c = (1209 + C57) % 10000;
    boolean flag = false;
    string label = "rec1209";
};
type R1210 record {
    int id;
    int a = C58 + 1210;
    int b = C47 * 7;
    int c = (1210 + C58) % 10000;
    boolean flag = true;
    string label = "rec1210";
};
type R1211 record {
    int id;
    int a = C59 + 1211;
    int b = C50 * 1;
    int c = (1211 + C59) % 10000;
    boolean flag = false;
    string label = "rec1211";
};
type R1212 record {
    int id;
    int a = C60 + 1212;
    int b = C53 * 2;
    int c = (1212 + C60) % 10000;
    boolean flag = true;
    string label = "rec1212";
};
type R1213 record {
    int id;
    int a = C61 + 1213;
    int b = C56 * 3;
    int c = (1213 + C61) % 10000;
    boolean flag = false;
    string label = "rec1213";
};
type R1214 record {
    int id;
    int a = C62 + 1214;
    int b = C59 * 4;
    int c = (1214 + C62) % 10000;
    boolean flag = true;
    string label = "rec1214";
};
type R1215 record {
    int id;
    int a = C63 + 1215;
    int b = C62 * 5;
    int c = (1215 + C63) % 10000;
    boolean flag = false;
    string label = "rec1215";
};
type R1216 record {
    int id;
    int a = C0 + 1216;
    int b = C1 * 6;
    int c = (1216 + C0) % 10000;
    boolean flag = true;
    string label = "rec1216";
};
type R1217 record {
    int id;
    int a = C1 + 1217;
    int b = C4 * 7;
    int c = (1217 + C1) % 10000;
    boolean flag = false;
    string label = "rec1217";
};
type R1218 record {
    int id;
    int a = C2 + 1218;
    int b = C7 * 1;
    int c = (1218 + C2) % 10000;
    boolean flag = true;
    string label = "rec1218";
};
type R1219 record {
    int id;
    int a = C3 + 1219;
    int b = C10 * 2;
    int c = (1219 + C3) % 10000;
    boolean flag = false;
    string label = "rec1219";
};
type R1220 record {
    int id;
    int a = C4 + 1220;
    int b = C13 * 3;
    int c = (1220 + C4) % 10000;
    boolean flag = true;
    string label = "rec1220";
};
type R1221 record {
    int id;
    int a = C5 + 1221;
    int b = C16 * 4;
    int c = (1221 + C5) % 10000;
    boolean flag = false;
    string label = "rec1221";
};
type R1222 record {
    int id;
    int a = C6 + 1222;
    int b = C19 * 5;
    int c = (1222 + C6) % 10000;
    boolean flag = true;
    string label = "rec1222";
};
type R1223 record {
    int id;
    int a = C7 + 1223;
    int b = C22 * 6;
    int c = (1223 + C7) % 10000;
    boolean flag = false;
    string label = "rec1223";
};
type R1224 record {
    int id;
    int a = C8 + 1224;
    int b = C25 * 7;
    int c = (1224 + C8) % 10000;
    boolean flag = true;
    string label = "rec1224";
};
type R1225 record {
    int id;
    int a = C9 + 1225;
    int b = C28 * 1;
    int c = (1225 + C9) % 10000;
    boolean flag = false;
    string label = "rec1225";
};
type R1226 record {
    int id;
    int a = C10 + 1226;
    int b = C31 * 2;
    int c = (1226 + C10) % 10000;
    boolean flag = true;
    string label = "rec1226";
};
type R1227 record {
    int id;
    int a = C11 + 1227;
    int b = C34 * 3;
    int c = (1227 + C11) % 10000;
    boolean flag = false;
    string label = "rec1227";
};
type R1228 record {
    int id;
    int a = C12 + 1228;
    int b = C37 * 4;
    int c = (1228 + C12) % 10000;
    boolean flag = true;
    string label = "rec1228";
};
type R1229 record {
    int id;
    int a = C13 + 1229;
    int b = C40 * 5;
    int c = (1229 + C13) % 10000;
    boolean flag = false;
    string label = "rec1229";
};
type R1230 record {
    int id;
    int a = C14 + 1230;
    int b = C43 * 6;
    int c = (1230 + C14) % 10000;
    boolean flag = true;
    string label = "rec1230";
};
type R1231 record {
    int id;
    int a = C15 + 1231;
    int b = C46 * 7;
    int c = (1231 + C15) % 10000;
    boolean flag = false;
    string label = "rec1231";
};
type R1232 record {
    int id;
    int a = C16 + 1232;
    int b = C49 * 1;
    int c = (1232 + C16) % 10000;
    boolean flag = true;
    string label = "rec1232";
};
type R1233 record {
    int id;
    int a = C17 + 1233;
    int b = C52 * 2;
    int c = (1233 + C17) % 10000;
    boolean flag = false;
    string label = "rec1233";
};
type R1234 record {
    int id;
    int a = C18 + 1234;
    int b = C55 * 3;
    int c = (1234 + C18) % 10000;
    boolean flag = true;
    string label = "rec1234";
};
type R1235 record {
    int id;
    int a = C19 + 1235;
    int b = C58 * 4;
    int c = (1235 + C19) % 10000;
    boolean flag = false;
    string label = "rec1235";
};
type R1236 record {
    int id;
    int a = C20 + 1236;
    int b = C61 * 5;
    int c = (1236 + C20) % 10000;
    boolean flag = true;
    string label = "rec1236";
};
type R1237 record {
    int id;
    int a = C21 + 1237;
    int b = C0 * 6;
    int c = (1237 + C21) % 10000;
    boolean flag = false;
    string label = "rec1237";
};
type R1238 record {
    int id;
    int a = C22 + 1238;
    int b = C3 * 7;
    int c = (1238 + C22) % 10000;
    boolean flag = true;
    string label = "rec1238";
};
type R1239 record {
    int id;
    int a = C23 + 1239;
    int b = C6 * 1;
    int c = (1239 + C23) % 10000;
    boolean flag = false;
    string label = "rec1239";
};
type R1240 record {
    int id;
    int a = C24 + 1240;
    int b = C9 * 2;
    int c = (1240 + C24) % 10000;
    boolean flag = true;
    string label = "rec1240";
};
type R1241 record {
    int id;
    int a = C25 + 1241;
    int b = C12 * 3;
    int c = (1241 + C25) % 10000;
    boolean flag = false;
    string label = "rec1241";
};
type R1242 record {
    int id;
    int a = C26 + 1242;
    int b = C15 * 4;
    int c = (1242 + C26) % 10000;
    boolean flag = true;
    string label = "rec1242";
};
type R1243 record {
    int id;
    int a = C27 + 1243;
    int b = C18 * 5;
    int c = (1243 + C27) % 10000;
    boolean flag = false;
    string label = "rec1243";
};
type R1244 record {
    int id;
    int a = C28 + 1244;
    int b = C21 * 6;
    int c = (1244 + C28) % 10000;
    boolean flag = true;
    string label = "rec1244";
};
type R1245 record {
    int id;
    int a = C29 + 1245;
    int b = C24 * 7;
    int c = (1245 + C29) % 10000;
    boolean flag = false;
    string label = "rec1245";
};
type R1246 record {
    int id;
    int a = C30 + 1246;
    int b = C27 * 1;
    int c = (1246 + C30) % 10000;
    boolean flag = true;
    string label = "rec1246";
};
type R1247 record {
    int id;
    int a = C31 + 1247;
    int b = C30 * 2;
    int c = (1247 + C31) % 10000;
    boolean flag = false;
    string label = "rec1247";
};
type R1248 record {
    int id;
    int a = C32 + 1248;
    int b = C33 * 3;
    int c = (1248 + C32) % 10000;
    boolean flag = true;
    string label = "rec1248";
};
type R1249 record {
    int id;
    int a = C33 + 1249;
    int b = C36 * 4;
    int c = (1249 + C33) % 10000;
    boolean flag = false;
    string label = "rec1249";
};
type R1250 record {
    int id;
    int a = C34 + 1250;
    int b = C39 * 5;
    int c = (1250 + C34) % 10000;
    boolean flag = true;
    string label = "rec1250";
};
type R1251 record {
    int id;
    int a = C35 + 1251;
    int b = C42 * 6;
    int c = (1251 + C35) % 10000;
    boolean flag = false;
    string label = "rec1251";
};
type R1252 record {
    int id;
    int a = C36 + 1252;
    int b = C45 * 7;
    int c = (1252 + C36) % 10000;
    boolean flag = true;
    string label = "rec1252";
};
type R1253 record {
    int id;
    int a = C37 + 1253;
    int b = C48 * 1;
    int c = (1253 + C37) % 10000;
    boolean flag = false;
    string label = "rec1253";
};
type R1254 record {
    int id;
    int a = C38 + 1254;
    int b = C51 * 2;
    int c = (1254 + C38) % 10000;
    boolean flag = true;
    string label = "rec1254";
};
type R1255 record {
    int id;
    int a = C39 + 1255;
    int b = C54 * 3;
    int c = (1255 + C39) % 10000;
    boolean flag = false;
    string label = "rec1255";
};
type R1256 record {
    int id;
    int a = C40 + 1256;
    int b = C57 * 4;
    int c = (1256 + C40) % 10000;
    boolean flag = true;
    string label = "rec1256";
};
type R1257 record {
    int id;
    int a = C41 + 1257;
    int b = C60 * 5;
    int c = (1257 + C41) % 10000;
    boolean flag = false;
    string label = "rec1257";
};
type R1258 record {
    int id;
    int a = C42 + 1258;
    int b = C63 * 6;
    int c = (1258 + C42) % 10000;
    boolean flag = true;
    string label = "rec1258";
};
type R1259 record {
    int id;
    int a = C43 + 1259;
    int b = C2 * 7;
    int c = (1259 + C43) % 10000;
    boolean flag = false;
    string label = "rec1259";
};
type R1260 record {
    int id;
    int a = C44 + 1260;
    int b = C5 * 1;
    int c = (1260 + C44) % 10000;
    boolean flag = true;
    string label = "rec1260";
};
type R1261 record {
    int id;
    int a = C45 + 1261;
    int b = C8 * 2;
    int c = (1261 + C45) % 10000;
    boolean flag = false;
    string label = "rec1261";
};
type R1262 record {
    int id;
    int a = C46 + 1262;
    int b = C11 * 3;
    int c = (1262 + C46) % 10000;
    boolean flag = true;
    string label = "rec1262";
};
type R1263 record {
    int id;
    int a = C47 + 1263;
    int b = C14 * 4;
    int c = (1263 + C47) % 10000;
    boolean flag = false;
    string label = "rec1263";
};
type R1264 record {
    int id;
    int a = C48 + 1264;
    int b = C17 * 5;
    int c = (1264 + C48) % 10000;
    boolean flag = true;
    string label = "rec1264";
};
type R1265 record {
    int id;
    int a = C49 + 1265;
    int b = C20 * 6;
    int c = (1265 + C49) % 10000;
    boolean flag = false;
    string label = "rec1265";
};
type R1266 record {
    int id;
    int a = C50 + 1266;
    int b = C23 * 7;
    int c = (1266 + C50) % 10000;
    boolean flag = true;
    string label = "rec1266";
};
type R1267 record {
    int id;
    int a = C51 + 1267;
    int b = C26 * 1;
    int c = (1267 + C51) % 10000;
    boolean flag = false;
    string label = "rec1267";
};
type R1268 record {
    int id;
    int a = C52 + 1268;
    int b = C29 * 2;
    int c = (1268 + C52) % 10000;
    boolean flag = true;
    string label = "rec1268";
};
type R1269 record {
    int id;
    int a = C53 + 1269;
    int b = C32 * 3;
    int c = (1269 + C53) % 10000;
    boolean flag = false;
    string label = "rec1269";
};
type R1270 record {
    int id;
    int a = C54 + 1270;
    int b = C35 * 4;
    int c = (1270 + C54) % 10000;
    boolean flag = true;
    string label = "rec1270";
};
type R1271 record {
    int id;
    int a = C55 + 1271;
    int b = C38 * 5;
    int c = (1271 + C55) % 10000;
    boolean flag = false;
    string label = "rec1271";
};
type R1272 record {
    int id;
    int a = C56 + 1272;
    int b = C41 * 6;
    int c = (1272 + C56) % 10000;
    boolean flag = true;
    string label = "rec1272";
};
type R1273 record {
    int id;
    int a = C57 + 1273;
    int b = C44 * 7;
    int c = (1273 + C57) % 10000;
    boolean flag = false;
    string label = "rec1273";
};
type R1274 record {
    int id;
    int a = C58 + 1274;
    int b = C47 * 1;
    int c = (1274 + C58) % 10000;
    boolean flag = true;
    string label = "rec1274";
};
type R1275 record {
    int id;
    int a = C59 + 1275;
    int b = C50 * 2;
    int c = (1275 + C59) % 10000;
    boolean flag = false;
    string label = "rec1275";
};
type R1276 record {
    int id;
    int a = C60 + 1276;
    int b = C53 * 3;
    int c = (1276 + C60) % 10000;
    boolean flag = true;
    string label = "rec1276";
};
type R1277 record {
    int id;
    int a = C61 + 1277;
    int b = C56 * 4;
    int c = (1277 + C61) % 10000;
    boolean flag = false;
    string label = "rec1277";
};
type R1278 record {
    int id;
    int a = C62 + 1278;
    int b = C59 * 5;
    int c = (1278 + C62) % 10000;
    boolean flag = true;
    string label = "rec1278";
};
type R1279 record {
    int id;
    int a = C63 + 1279;
    int b = C62 * 6;
    int c = (1279 + C63) % 10000;
    boolean flag = false;
    string label = "rec1279";
};
type R1280 record {
    int id;
    int a = C0 + 1280;
    int b = C1 * 7;
    int c = (1280 + C0) % 10000;
    boolean flag = true;
    string label = "rec1280";
};
type R1281 record {
    int id;
    int a = C1 + 1281;
    int b = C4 * 1;
    int c = (1281 + C1) % 10000;
    boolean flag = false;
    string label = "rec1281";
};
type R1282 record {
    int id;
    int a = C2 + 1282;
    int b = C7 * 2;
    int c = (1282 + C2) % 10000;
    boolean flag = true;
    string label = "rec1282";
};
type R1283 record {
    int id;
    int a = C3 + 1283;
    int b = C10 * 3;
    int c = (1283 + C3) % 10000;
    boolean flag = false;
    string label = "rec1283";
};
type R1284 record {
    int id;
    int a = C4 + 1284;
    int b = C13 * 4;
    int c = (1284 + C4) % 10000;
    boolean flag = true;
    string label = "rec1284";
};
type R1285 record {
    int id;
    int a = C5 + 1285;
    int b = C16 * 5;
    int c = (1285 + C5) % 10000;
    boolean flag = false;
    string label = "rec1285";
};
type R1286 record {
    int id;
    int a = C6 + 1286;
    int b = C19 * 6;
    int c = (1286 + C6) % 10000;
    boolean flag = true;
    string label = "rec1286";
};
type R1287 record {
    int id;
    int a = C7 + 1287;
    int b = C22 * 7;
    int c = (1287 + C7) % 10000;
    boolean flag = false;
    string label = "rec1287";
};
type R1288 record {
    int id;
    int a = C8 + 1288;
    int b = C25 * 1;
    int c = (1288 + C8) % 10000;
    boolean flag = true;
    string label = "rec1288";
};
type R1289 record {
    int id;
    int a = C9 + 1289;
    int b = C28 * 2;
    int c = (1289 + C9) % 10000;
    boolean flag = false;
    string label = "rec1289";
};
type R1290 record {
    int id;
    int a = C10 + 1290;
    int b = C31 * 3;
    int c = (1290 + C10) % 10000;
    boolean flag = true;
    string label = "rec1290";
};
type R1291 record {
    int id;
    int a = C11 + 1291;
    int b = C34 * 4;
    int c = (1291 + C11) % 10000;
    boolean flag = false;
    string label = "rec1291";
};
type R1292 record {
    int id;
    int a = C12 + 1292;
    int b = C37 * 5;
    int c = (1292 + C12) % 10000;
    boolean flag = true;
    string label = "rec1292";
};
type R1293 record {
    int id;
    int a = C13 + 1293;
    int b = C40 * 6;
    int c = (1293 + C13) % 10000;
    boolean flag = false;
    string label = "rec1293";
};
type R1294 record {
    int id;
    int a = C14 + 1294;
    int b = C43 * 7;
    int c = (1294 + C14) % 10000;
    boolean flag = true;
    string label = "rec1294";
};
type R1295 record {
    int id;
    int a = C15 + 1295;
    int b = C46 * 1;
    int c = (1295 + C15) % 10000;
    boolean flag = false;
    string label = "rec1295";
};
type R1296 record {
    int id;
    int a = C16 + 1296;
    int b = C49 * 2;
    int c = (1296 + C16) % 10000;
    boolean flag = true;
    string label = "rec1296";
};
type R1297 record {
    int id;
    int a = C17 + 1297;
    int b = C52 * 3;
    int c = (1297 + C17) % 10000;
    boolean flag = false;
    string label = "rec1297";
};
type R1298 record {
    int id;
    int a = C18 + 1298;
    int b = C55 * 4;
    int c = (1298 + C18) % 10000;
    boolean flag = true;
    string label = "rec1298";
};
type R1299 record {
    int id;
    int a = C19 + 1299;
    int b = C58 * 5;
    int c = (1299 + C19) % 10000;
    boolean flag = false;
    string label = "rec1299";
};
type R1300 record {
    int id;
    int a = C20 + 1300;
    int b = C61 * 6;
    int c = (1300 + C20) % 10000;
    boolean flag = true;
    string label = "rec1300";
};
type R1301 record {
    int id;
    int a = C21 + 1301;
    int b = C0 * 7;
    int c = (1301 + C21) % 10000;
    boolean flag = false;
    string label = "rec1301";
};
type R1302 record {
    int id;
    int a = C22 + 1302;
    int b = C3 * 1;
    int c = (1302 + C22) % 10000;
    boolean flag = true;
    string label = "rec1302";
};
type R1303 record {
    int id;
    int a = C23 + 1303;
    int b = C6 * 2;
    int c = (1303 + C23) % 10000;
    boolean flag = false;
    string label = "rec1303";
};
type R1304 record {
    int id;
    int a = C24 + 1304;
    int b = C9 * 3;
    int c = (1304 + C24) % 10000;
    boolean flag = true;
    string label = "rec1304";
};
type R1305 record {
    int id;
    int a = C25 + 1305;
    int b = C12 * 4;
    int c = (1305 + C25) % 10000;
    boolean flag = false;
    string label = "rec1305";
};
type R1306 record {
    int id;
    int a = C26 + 1306;
    int b = C15 * 5;
    int c = (1306 + C26) % 10000;
    boolean flag = true;
    string label = "rec1306";
};
type R1307 record {
    int id;
    int a = C27 + 1307;
    int b = C18 * 6;
    int c = (1307 + C27) % 10000;
    boolean flag = false;
    string label = "rec1307";
};
type R1308 record {
    int id;
    int a = C28 + 1308;
    int b = C21 * 7;
    int c = (1308 + C28) % 10000;
    boolean flag = true;
    string label = "rec1308";
};
type R1309 record {
    int id;
    int a = C29 + 1309;
    int b = C24 * 1;
    int c = (1309 + C29) % 10000;
    boolean flag = false;
    string label = "rec1309";
};
type R1310 record {
    int id;
    int a = C30 + 1310;
    int b = C27 * 2;
    int c = (1310 + C30) % 10000;
    boolean flag = true;
    string label = "rec1310";
};
type R1311 record {
    int id;
    int a = C31 + 1311;
    int b = C30 * 3;
    int c = (1311 + C31) % 10000;
    boolean flag = false;
    string label = "rec1311";
};
type R1312 record {
    int id;
    int a = C32 + 1312;
    int b = C33 * 4;
    int c = (1312 + C32) % 10000;
    boolean flag = true;
    string label = "rec1312";
};
type R1313 record {
    int id;
    int a = C33 + 1313;
    int b = C36 * 5;
    int c = (1313 + C33) % 10000;
    boolean flag = false;
    string label = "rec1313";
};
type R1314 record {
    int id;
    int a = C34 + 1314;
    int b = C39 * 6;
    int c = (1314 + C34) % 10000;
    boolean flag = true;
    string label = "rec1314";
};
type R1315 record {
    int id;
    int a = C35 + 1315;
    int b = C42 * 7;
    int c = (1315 + C35) % 10000;
    boolean flag = false;
    string label = "rec1315";
};
type R1316 record {
    int id;
    int a = C36 + 1316;
    int b = C45 * 1;
    int c = (1316 + C36) % 10000;
    boolean flag = true;
    string label = "rec1316";
};
type R1317 record {
    int id;
    int a = C37 + 1317;
    int b = C48 * 2;
    int c = (1317 + C37) % 10000;
    boolean flag = false;
    string label = "rec1317";
};
type R1318 record {
    int id;
    int a = C38 + 1318;
    int b = C51 * 3;
    int c = (1318 + C38) % 10000;
    boolean flag = true;
    string label = "rec1318";
};
type R1319 record {
    int id;
    int a = C39 + 1319;
    int b = C54 * 4;
    int c = (1319 + C39) % 10000;
    boolean flag = false;
    string label = "rec1319";
};
type R1320 record {
    int id;
    int a = C40 + 1320;
    int b = C57 * 5;
    int c = (1320 + C40) % 10000;
    boolean flag = true;
    string label = "rec1320";
};
type R1321 record {
    int id;
    int a = C41 + 1321;
    int b = C60 * 6;
    int c = (1321 + C41) % 10000;
    boolean flag = false;
    string label = "rec1321";
};
type R1322 record {
    int id;
    int a = C42 + 1322;
    int b = C63 * 7;
    int c = (1322 + C42) % 10000;
    boolean flag = true;
    string label = "rec1322";
};
type R1323 record {
    int id;
    int a = C43 + 1323;
    int b = C2 * 1;
    int c = (1323 + C43) % 10000;
    boolean flag = false;
    string label = "rec1323";
};
type R1324 record {
    int id;
    int a = C44 + 1324;
    int b = C5 * 2;
    int c = (1324 + C44) % 10000;
    boolean flag = true;
    string label = "rec1324";
};
type R1325 record {
    int id;
    int a = C45 + 1325;
    int b = C8 * 3;
    int c = (1325 + C45) % 10000;
    boolean flag = false;
    string label = "rec1325";
};
type R1326 record {
    int id;
    int a = C46 + 1326;
    int b = C11 * 4;
    int c = (1326 + C46) % 10000;
    boolean flag = true;
    string label = "rec1326";
};
type R1327 record {
    int id;
    int a = C47 + 1327;
    int b = C14 * 5;
    int c = (1327 + C47) % 10000;
    boolean flag = false;
    string label = "rec1327";
};
type R1328 record {
    int id;
    int a = C48 + 1328;
    int b = C17 * 6;
    int c = (1328 + C48) % 10000;
    boolean flag = true;
    string label = "rec1328";
};
type R1329 record {
    int id;
    int a = C49 + 1329;
    int b = C20 * 7;
    int c = (1329 + C49) % 10000;
    boolean flag = false;
    string label = "rec1329";
};
type R1330 record {
    int id;
    int a = C50 + 1330;
    int b = C23 * 1;
    int c = (1330 + C50) % 10000;
    boolean flag = true;
    string label = "rec1330";
};
type R1331 record {
    int id;
    int a = C51 + 1331;
    int b = C26 * 2;
    int c = (1331 + C51) % 10000;
    boolean flag = false;
    string label = "rec1331";
};
type R1332 record {
    int id;
    int a = C52 + 1332;
    int b = C29 * 3;
    int c = (1332 + C52) % 10000;
    boolean flag = true;
    string label = "rec1332";
};
type R1333 record {
    int id;
    int a = C53 + 1333;
    int b = C32 * 4;
    int c = (1333 + C53) % 10000;
    boolean flag = false;
    string label = "rec1333";
};
type R1334 record {
    int id;
    int a = C54 + 1334;
    int b = C35 * 5;
    int c = (1334 + C54) % 10000;
    boolean flag = true;
    string label = "rec1334";
};
type R1335 record {
    int id;
    int a = C55 + 1335;
    int b = C38 * 6;
    int c = (1335 + C55) % 10000;
    boolean flag = false;
    string label = "rec1335";
};
type R1336 record {
    int id;
    int a = C56 + 1336;
    int b = C41 * 7;
    int c = (1336 + C56) % 10000;
    boolean flag = true;
    string label = "rec1336";
};
type R1337 record {
    int id;
    int a = C57 + 1337;
    int b = C44 * 1;
    int c = (1337 + C57) % 10000;
    boolean flag = false;
    string label = "rec1337";
};
type R1338 record {
    int id;
    int a = C58 + 1338;
    int b = C47 * 2;
    int c = (1338 + C58) % 10000;
    boolean flag = true;
    string label = "rec1338";
};
type R1339 record {
    int id;
    int a = C59 + 1339;
    int b = C50 * 3;
    int c = (1339 + C59) % 10000;
    boolean flag = false;
    string label = "rec1339";
};
type R1340 record {
    int id;
    int a = C60 + 1340;
    int b = C53 * 4;
    int c = (1340 + C60) % 10000;
    boolean flag = true;
    string label = "rec1340";
};
type R1341 record {
    int id;
    int a = C61 + 1341;
    int b = C56 * 5;
    int c = (1341 + C61) % 10000;
    boolean flag = false;
    string label = "rec1341";
};
type R1342 record {
    int id;
    int a = C62 + 1342;
    int b = C59 * 6;
    int c = (1342 + C62) % 10000;
    boolean flag = true;
    string label = "rec1342";
};
type R1343 record {
    int id;
    int a = C63 + 1343;
    int b = C62 * 7;
    int c = (1343 + C63) % 10000;
    boolean flag = false;
    string label = "rec1343";
};
type R1344 record {
    int id;
    int a = C0 + 1344;
    int b = C1 * 1;
    int c = (1344 + C0) % 10000;
    boolean flag = true;
    string label = "rec1344";
};
type R1345 record {
    int id;
    int a = C1 + 1345;
    int b = C4 * 2;
    int c = (1345 + C1) % 10000;
    boolean flag = false;
    string label = "rec1345";
};
type R1346 record {
    int id;
    int a = C2 + 1346;
    int b = C7 * 3;
    int c = (1346 + C2) % 10000;
    boolean flag = true;
    string label = "rec1346";
};
type R1347 record {
    int id;
    int a = C3 + 1347;
    int b = C10 * 4;
    int c = (1347 + C3) % 10000;
    boolean flag = false;
    string label = "rec1347";
};
type R1348 record {
    int id;
    int a = C4 + 1348;
    int b = C13 * 5;
    int c = (1348 + C4) % 10000;
    boolean flag = true;
    string label = "rec1348";
};
type R1349 record {
    int id;
    int a = C5 + 1349;
    int b = C16 * 6;
    int c = (1349 + C5) % 10000;
    boolean flag = false;
    string label = "rec1349";
};
type R1350 record {
    int id;
    int a = C6 + 1350;
    int b = C19 * 7;
    int c = (1350 + C6) % 10000;
    boolean flag = true;
    string label = "rec1350";
};
type R1351 record {
    int id;
    int a = C7 + 1351;
    int b = C22 * 1;
    int c = (1351 + C7) % 10000;
    boolean flag = false;
    string label = "rec1351";
};
type R1352 record {
    int id;
    int a = C8 + 1352;
    int b = C25 * 2;
    int c = (1352 + C8) % 10000;
    boolean flag = true;
    string label = "rec1352";
};
type R1353 record {
    int id;
    int a = C9 + 1353;
    int b = C28 * 3;
    int c = (1353 + C9) % 10000;
    boolean flag = false;
    string label = "rec1353";
};
type R1354 record {
    int id;
    int a = C10 + 1354;
    int b = C31 * 4;
    int c = (1354 + C10) % 10000;
    boolean flag = true;
    string label = "rec1354";
};
type R1355 record {
    int id;
    int a = C11 + 1355;
    int b = C34 * 5;
    int c = (1355 + C11) % 10000;
    boolean flag = false;
    string label = "rec1355";
};
type R1356 record {
    int id;
    int a = C12 + 1356;
    int b = C37 * 6;
    int c = (1356 + C12) % 10000;
    boolean flag = true;
    string label = "rec1356";
};
type R1357 record {
    int id;
    int a = C13 + 1357;
    int b = C40 * 7;
    int c = (1357 + C13) % 10000;
    boolean flag = false;
    string label = "rec1357";
};
type R1358 record {
    int id;
    int a = C14 + 1358;
    int b = C43 * 1;
    int c = (1358 + C14) % 10000;
    boolean flag = true;
    string label = "rec1358";
};
type R1359 record {
    int id;
    int a = C15 + 1359;
    int b = C46 * 2;
    int c = (1359 + C15) % 10000;
    boolean flag = false;
    string label = "rec1359";
};
type R1360 record {
    int id;
    int a = C16 + 1360;
    int b = C49 * 3;
    int c = (1360 + C16) % 10000;
    boolean flag = true;
    string label = "rec1360";
};
type R1361 record {
    int id;
    int a = C17 + 1361;
    int b = C52 * 4;
    int c = (1361 + C17) % 10000;
    boolean flag = false;
    string label = "rec1361";
};
type R1362 record {
    int id;
    int a = C18 + 1362;
    int b = C55 * 5;
    int c = (1362 + C18) % 10000;
    boolean flag = true;
    string label = "rec1362";
};
type R1363 record {
    int id;
    int a = C19 + 1363;
    int b = C58 * 6;
    int c = (1363 + C19) % 10000;
    boolean flag = false;
    string label = "rec1363";
};
type R1364 record {
    int id;
    int a = C20 + 1364;
    int b = C61 * 7;
    int c = (1364 + C20) % 10000;
    boolean flag = true;
    string label = "rec1364";
};
type R1365 record {
    int id;
    int a = C21 + 1365;
    int b = C0 * 1;
    int c = (1365 + C21) % 10000;
    boolean flag = false;
    string label = "rec1365";
};
type R1366 record {
    int id;
    int a = C22 + 1366;
    int b = C3 * 2;
    int c = (1366 + C22) % 10000;
    boolean flag = true;
    string label = "rec1366";
};
type R1367 record {
    int id;
    int a = C23 + 1367;
    int b = C6 * 3;
    int c = (1367 + C23) % 10000;
    boolean flag = false;
    string label = "rec1367";
};
type R1368 record {
    int id;
    int a = C24 + 1368;
    int b = C9 * 4;
    int c = (1368 + C24) % 10000;
    boolean flag = true;
    string label = "rec1368";
};
type R1369 record {
    int id;
    int a = C25 + 1369;
    int b = C12 * 5;
    int c = (1369 + C25) % 10000;
    boolean flag = false;
    string label = "rec1369";
};
type R1370 record {
    int id;
    int a = C26 + 1370;
    int b = C15 * 6;
    int c = (1370 + C26) % 10000;
    boolean flag = true;
    string label = "rec1370";
};
type R1371 record {
    int id;
    int a = C27 + 1371;
    int b = C18 * 7;
    int c = (1371 + C27) % 10000;
    boolean flag = false;
    string label = "rec1371";
};
type R1372 record {
    int id;
    int a = C28 + 1372;
    int b = C21 * 1;
    int c = (1372 + C28) % 10000;
    boolean flag = true;
    string label = "rec1372";
};
type R1373 record {
    int id;
    int a = C29 + 1373;
    int b = C24 * 2;
    int c = (1373 + C29) % 10000;
    boolean flag = false;
    string label = "rec1373";
};
type R1374 record {
    int id;
    int a = C30 + 1374;
    int b = C27 * 3;
    int c = (1374 + C30) % 10000;
    boolean flag = true;
    string label = "rec1374";
};
type R1375 record {
    int id;
    int a = C31 + 1375;
    int b = C30 * 4;
    int c = (1375 + C31) % 10000;
    boolean flag = false;
    string label = "rec1375";
};
type R1376 record {
    int id;
    int a = C32 + 1376;
    int b = C33 * 5;
    int c = (1376 + C32) % 10000;
    boolean flag = true;
    string label = "rec1376";
};
type R1377 record {
    int id;
    int a = C33 + 1377;
    int b = C36 * 6;
    int c = (1377 + C33) % 10000;
    boolean flag = false;
    string label = "rec1377";
};
type R1378 record {
    int id;
    int a = C34 + 1378;
    int b = C39 * 7;
    int c = (1378 + C34) % 10000;
    boolean flag = true;
    string label = "rec1378";
};
type R1379 record {
    int id;
    int a = C35 + 1379;
    int b = C42 * 1;
    int c = (1379 + C35) % 10000;
    boolean flag = false;
    string label = "rec1379";
};
type R1380 record {
    int id;
    int a = C36 + 1380;
    int b = C45 * 2;
    int c = (1380 + C36) % 10000;
    boolean flag = true;
    string label = "rec1380";
};
type R1381 record {
    int id;
    int a = C37 + 1381;
    int b = C48 * 3;
    int c = (1381 + C37) % 10000;
    boolean flag = false;
    string label = "rec1381";
};
type R1382 record {
    int id;
    int a = C38 + 1382;
    int b = C51 * 4;
    int c = (1382 + C38) % 10000;
    boolean flag = true;
    string label = "rec1382";
};
type R1383 record {
    int id;
    int a = C39 + 1383;
    int b = C54 * 5;
    int c = (1383 + C39) % 10000;
    boolean flag = false;
    string label = "rec1383";
};
type R1384 record {
    int id;
    int a = C40 + 1384;
    int b = C57 * 6;
    int c = (1384 + C40) % 10000;
    boolean flag = true;
    string label = "rec1384";
};
type R1385 record {
    int id;
    int a = C41 + 1385;
    int b = C60 * 7;
    int c = (1385 + C41) % 10000;
    boolean flag = false;
    string label = "rec1385";
};
type R1386 record {
    int id;
    int a = C42 + 1386;
    int b = C63 * 1;
    int c = (1386 + C42) % 10000;
    boolean flag = true;
    string label = "rec1386";
};
type R1387 record {
    int id;
    int a = C43 + 1387;
    int b = C2 * 2;
    int c = (1387 + C43) % 10000;
    boolean flag = false;
    string label = "rec1387";
};
type R1388 record {
    int id;
    int a = C44 + 1388;
    int b = C5 * 3;
    int c = (1388 + C44) % 10000;
    boolean flag = true;
    string label = "rec1388";
};
type R1389 record {
    int id;
    int a = C45 + 1389;
    int b = C8 * 4;
    int c = (1389 + C45) % 10000;
    boolean flag = false;
    string label = "rec1389";
};
type R1390 record {
    int id;
    int a = C46 + 1390;
    int b = C11 * 5;
    int c = (1390 + C46) % 10000;
    boolean flag = true;
    string label = "rec1390";
};
type R1391 record {
    int id;
    int a = C47 + 1391;
    int b = C14 * 6;
    int c = (1391 + C47) % 10000;
    boolean flag = false;
    string label = "rec1391";
};
type R1392 record {
    int id;
    int a = C48 + 1392;
    int b = C17 * 7;
    int c = (1392 + C48) % 10000;
    boolean flag = true;
    string label = "rec1392";
};
type R1393 record {
    int id;
    int a = C49 + 1393;
    int b = C20 * 1;
    int c = (1393 + C49) % 10000;
    boolean flag = false;
    string label = "rec1393";
};
type R1394 record {
    int id;
    int a = C50 + 1394;
    int b = C23 * 2;
    int c = (1394 + C50) % 10000;
    boolean flag = true;
    string label = "rec1394";
};
type R1395 record {
    int id;
    int a = C51 + 1395;
    int b = C26 * 3;
    int c = (1395 + C51) % 10000;
    boolean flag = false;
    string label = "rec1395";
};
type R1396 record {
    int id;
    int a = C52 + 1396;
    int b = C29 * 4;
    int c = (1396 + C52) % 10000;
    boolean flag = true;
    string label = "rec1396";
};
type R1397 record {
    int id;
    int a = C53 + 1397;
    int b = C32 * 5;
    int c = (1397 + C53) % 10000;
    boolean flag = false;
    string label = "rec1397";
};
type R1398 record {
    int id;
    int a = C54 + 1398;
    int b = C35 * 6;
    int c = (1398 + C54) % 10000;
    boolean flag = true;
    string label = "rec1398";
};
type R1399 record {
    int id;
    int a = C55 + 1399;
    int b = C38 * 7;
    int c = (1399 + C55) % 10000;
    boolean flag = false;
    string label = "rec1399";
};
type R1400 record {
    int id;
    int a = C56 + 1400;
    int b = C41 * 1;
    int c = (1400 + C56) % 10000;
    boolean flag = true;
    string label = "rec1400";
};
type R1401 record {
    int id;
    int a = C57 + 1401;
    int b = C44 * 2;
    int c = (1401 + C57) % 10000;
    boolean flag = false;
    string label = "rec1401";
};
type R1402 record {
    int id;
    int a = C58 + 1402;
    int b = C47 * 3;
    int c = (1402 + C58) % 10000;
    boolean flag = true;
    string label = "rec1402";
};
type R1403 record {
    int id;
    int a = C59 + 1403;
    int b = C50 * 4;
    int c = (1403 + C59) % 10000;
    boolean flag = false;
    string label = "rec1403";
};
type R1404 record {
    int id;
    int a = C60 + 1404;
    int b = C53 * 5;
    int c = (1404 + C60) % 10000;
    boolean flag = true;
    string label = "rec1404";
};
type R1405 record {
    int id;
    int a = C61 + 1405;
    int b = C56 * 6;
    int c = (1405 + C61) % 10000;
    boolean flag = false;
    string label = "rec1405";
};
type R1406 record {
    int id;
    int a = C62 + 1406;
    int b = C59 * 7;
    int c = (1406 + C62) % 10000;
    boolean flag = true;
    string label = "rec1406";
};
type R1407 record {
    int id;
    int a = C63 + 1407;
    int b = C62 * 1;
    int c = (1407 + C63) % 10000;
    boolean flag = false;
    string label = "rec1407";
};
type R1408 record {
    int id;
    int a = C0 + 1408;
    int b = C1 * 2;
    int c = (1408 + C0) % 10000;
    boolean flag = true;
    string label = "rec1408";
};
type R1409 record {
    int id;
    int a = C1 + 1409;
    int b = C4 * 3;
    int c = (1409 + C1) % 10000;
    boolean flag = false;
    string label = "rec1409";
};
type R1410 record {
    int id;
    int a = C2 + 1410;
    int b = C7 * 4;
    int c = (1410 + C2) % 10000;
    boolean flag = true;
    string label = "rec1410";
};
type R1411 record {
    int id;
    int a = C3 + 1411;
    int b = C10 * 5;
    int c = (1411 + C3) % 10000;
    boolean flag = false;
    string label = "rec1411";
};
type R1412 record {
    int id;
    int a = C4 + 1412;
    int b = C13 * 6;
    int c = (1412 + C4) % 10000;
    boolean flag = true;
    string label = "rec1412";
};
type R1413 record {
    int id;
    int a = C5 + 1413;
    int b = C16 * 7;
    int c = (1413 + C5) % 10000;
    boolean flag = false;
    string label = "rec1413";
};
type R1414 record {
    int id;
    int a = C6 + 1414;
    int b = C19 * 1;
    int c = (1414 + C6) % 10000;
    boolean flag = true;
    string label = "rec1414";
};
type R1415 record {
    int id;
    int a = C7 + 1415;
    int b = C22 * 2;
    int c = (1415 + C7) % 10000;
    boolean flag = false;
    string label = "rec1415";
};
type R1416 record {
    int id;
    int a = C8 + 1416;
    int b = C25 * 3;
    int c = (1416 + C8) % 10000;
    boolean flag = true;
    string label = "rec1416";
};
type R1417 record {
    int id;
    int a = C9 + 1417;
    int b = C28 * 4;
    int c = (1417 + C9) % 10000;
    boolean flag = false;
    string label = "rec1417";
};
type R1418 record {
    int id;
    int a = C10 + 1418;
    int b = C31 * 5;
    int c = (1418 + C10) % 10000;
    boolean flag = true;
    string label = "rec1418";
};
type R1419 record {
    int id;
    int a = C11 + 1419;
    int b = C34 * 6;
    int c = (1419 + C11) % 10000;
    boolean flag = false;
    string label = "rec1419";
};
type R1420 record {
    int id;
    int a = C12 + 1420;
    int b = C37 * 7;
    int c = (1420 + C12) % 10000;
    boolean flag = true;
    string label = "rec1420";
};
type R1421 record {
    int id;
    int a = C13 + 1421;
    int b = C40 * 1;
    int c = (1421 + C13) % 10000;
    boolean flag = false;
    string label = "rec1421";
};
type R1422 record {
    int id;
    int a = C14 + 1422;
    int b = C43 * 2;
    int c = (1422 + C14) % 10000;
    boolean flag = true;
    string label = "rec1422";
};
type R1423 record {
    int id;
    int a = C15 + 1423;
    int b = C46 * 3;
    int c = (1423 + C15) % 10000;
    boolean flag = false;
    string label = "rec1423";
};
type R1424 record {
    int id;
    int a = C16 + 1424;
    int b = C49 * 4;
    int c = (1424 + C16) % 10000;
    boolean flag = true;
    string label = "rec1424";
};
type R1425 record {
    int id;
    int a = C17 + 1425;
    int b = C52 * 5;
    int c = (1425 + C17) % 10000;
    boolean flag = false;
    string label = "rec1425";
};
type R1426 record {
    int id;
    int a = C18 + 1426;
    int b = C55 * 6;
    int c = (1426 + C18) % 10000;
    boolean flag = true;
    string label = "rec1426";
};
type R1427 record {
    int id;
    int a = C19 + 1427;
    int b = C58 * 7;
    int c = (1427 + C19) % 10000;
    boolean flag = false;
    string label = "rec1427";
};
type R1428 record {
    int id;
    int a = C20 + 1428;
    int b = C61 * 1;
    int c = (1428 + C20) % 10000;
    boolean flag = true;
    string label = "rec1428";
};
type R1429 record {
    int id;
    int a = C21 + 1429;
    int b = C0 * 2;
    int c = (1429 + C21) % 10000;
    boolean flag = false;
    string label = "rec1429";
};
type R1430 record {
    int id;
    int a = C22 + 1430;
    int b = C3 * 3;
    int c = (1430 + C22) % 10000;
    boolean flag = true;
    string label = "rec1430";
};
type R1431 record {
    int id;
    int a = C23 + 1431;
    int b = C6 * 4;
    int c = (1431 + C23) % 10000;
    boolean flag = false;
    string label = "rec1431";
};
type R1432 record {
    int id;
    int a = C24 + 1432;
    int b = C9 * 5;
    int c = (1432 + C24) % 10000;
    boolean flag = true;
    string label = "rec1432";
};
type R1433 record {
    int id;
    int a = C25 + 1433;
    int b = C12 * 6;
    int c = (1433 + C25) % 10000;
    boolean flag = false;
    string label = "rec1433";
};
type R1434 record {
    int id;
    int a = C26 + 1434;
    int b = C15 * 7;
    int c = (1434 + C26) % 10000;
    boolean flag = true;
    string label = "rec1434";
};
type R1435 record {
    int id;
    int a = C27 + 1435;
    int b = C18 * 1;
    int c = (1435 + C27) % 10000;
    boolean flag = false;
    string label = "rec1435";
};
type R1436 record {
    int id;
    int a = C28 + 1436;
    int b = C21 * 2;
    int c = (1436 + C28) % 10000;
    boolean flag = true;
    string label = "rec1436";
};
type R1437 record {
    int id;
    int a = C29 + 1437;
    int b = C24 * 3;
    int c = (1437 + C29) % 10000;
    boolean flag = false;
    string label = "rec1437";
};
type R1438 record {
    int id;
    int a = C30 + 1438;
    int b = C27 * 4;
    int c = (1438 + C30) % 10000;
    boolean flag = true;
    string label = "rec1438";
};
type R1439 record {
    int id;
    int a = C31 + 1439;
    int b = C30 * 5;
    int c = (1439 + C31) % 10000;
    boolean flag = false;
    string label = "rec1439";
};
type R1440 record {
    int id;
    int a = C32 + 1440;
    int b = C33 * 6;
    int c = (1440 + C32) % 10000;
    boolean flag = true;
    string label = "rec1440";
};
type R1441 record {
    int id;
    int a = C33 + 1441;
    int b = C36 * 7;
    int c = (1441 + C33) % 10000;
    boolean flag = false;
    string label = "rec1441";
};
type R1442 record {
    int id;
    int a = C34 + 1442;
    int b = C39 * 1;
    int c = (1442 + C34) % 10000;
    boolean flag = true;
    string label = "rec1442";
};
type R1443 record {
    int id;
    int a = C35 + 1443;
    int b = C42 * 2;
    int c = (1443 + C35) % 10000;
    boolean flag = false;
    string label = "rec1443";
};
type R1444 record {
    int id;
    int a = C36 + 1444;
    int b = C45 * 3;
    int c = (1444 + C36) % 10000;
    boolean flag = true;
    string label = "rec1444";
};
type R1445 record {
    int id;
    int a = C37 + 1445;
    int b = C48 * 4;
    int c = (1445 + C37) % 10000;
    boolean flag = false;
    string label = "rec1445";
};
type R1446 record {
    int id;
    int a = C38 + 1446;
    int b = C51 * 5;
    int c = (1446 + C38) % 10000;
    boolean flag = true;
    string label = "rec1446";
};
type R1447 record {
    int id;
    int a = C39 + 1447;
    int b = C54 * 6;
    int c = (1447 + C39) % 10000;
    boolean flag = false;
    string label = "rec1447";
};
type R1448 record {
    int id;
    int a = C40 + 1448;
    int b = C57 * 7;
    int c = (1448 + C40) % 10000;
    boolean flag = true;
    string label = "rec1448";
};
type R1449 record {
    int id;
    int a = C41 + 1449;
    int b = C60 * 1;
    int c = (1449 + C41) % 10000;
    boolean flag = false;
    string label = "rec1449";
};
type R1450 record {
    int id;
    int a = C42 + 1450;
    int b = C63 * 2;
    int c = (1450 + C42) % 10000;
    boolean flag = true;
    string label = "rec1450";
};
type R1451 record {
    int id;
    int a = C43 + 1451;
    int b = C2 * 3;
    int c = (1451 + C43) % 10000;
    boolean flag = false;
    string label = "rec1451";
};
type R1452 record {
    int id;
    int a = C44 + 1452;
    int b = C5 * 4;
    int c = (1452 + C44) % 10000;
    boolean flag = true;
    string label = "rec1452";
};
type R1453 record {
    int id;
    int a = C45 + 1453;
    int b = C8 * 5;
    int c = (1453 + C45) % 10000;
    boolean flag = false;
    string label = "rec1453";
};
type R1454 record {
    int id;
    int a = C46 + 1454;
    int b = C11 * 6;
    int c = (1454 + C46) % 10000;
    boolean flag = true;
    string label = "rec1454";
};
type R1455 record {
    int id;
    int a = C47 + 1455;
    int b = C14 * 7;
    int c = (1455 + C47) % 10000;
    boolean flag = false;
    string label = "rec1455";
};
type R1456 record {
    int id;
    int a = C48 + 1456;
    int b = C17 * 1;
    int c = (1456 + C48) % 10000;
    boolean flag = true;
    string label = "rec1456";
};
type R1457 record {
    int id;
    int a = C49 + 1457;
    int b = C20 * 2;
    int c = (1457 + C49) % 10000;
    boolean flag = false;
    string label = "rec1457";
};
type R1458 record {
    int id;
    int a = C50 + 1458;
    int b = C23 * 3;
    int c = (1458 + C50) % 10000;
    boolean flag = true;
    string label = "rec1458";
};
type R1459 record {
    int id;
    int a = C51 + 1459;
    int b = C26 * 4;
    int c = (1459 + C51) % 10000;
    boolean flag = false;
    string label = "rec1459";
};
type R1460 record {
    int id;
    int a = C52 + 1460;
    int b = C29 * 5;
    int c = (1460 + C52) % 10000;
    boolean flag = true;
    string label = "rec1460";
};
type R1461 record {
    int id;
    int a = C53 + 1461;
    int b = C32 * 6;
    int c = (1461 + C53) % 10000;
    boolean flag = false;
    string label = "rec1461";
};
type R1462 record {
    int id;
    int a = C54 + 1462;
    int b = C35 * 7;
    int c = (1462 + C54) % 10000;
    boolean flag = true;
    string label = "rec1462";
};
type R1463 record {
    int id;
    int a = C55 + 1463;
    int b = C38 * 1;
    int c = (1463 + C55) % 10000;
    boolean flag = false;
    string label = "rec1463";
};
type R1464 record {
    int id;
    int a = C56 + 1464;
    int b = C41 * 2;
    int c = (1464 + C56) % 10000;
    boolean flag = true;
    string label = "rec1464";
};
type R1465 record {
    int id;
    int a = C57 + 1465;
    int b = C44 * 3;
    int c = (1465 + C57) % 10000;
    boolean flag = false;
    string label = "rec1465";
};
type R1466 record {
    int id;
    int a = C58 + 1466;
    int b = C47 * 4;
    int c = (1466 + C58) % 10000;
    boolean flag = true;
    string label = "rec1466";
};
type R1467 record {
    int id;
    int a = C59 + 1467;
    int b = C50 * 5;
    int c = (1467 + C59) % 10000;
    boolean flag = false;
    string label = "rec1467";
};
type R1468 record {
    int id;
    int a = C60 + 1468;
    int b = C53 * 6;
    int c = (1468 + C60) % 10000;
    boolean flag = true;
    string label = "rec1468";
};
type R1469 record {
    int id;
    int a = C61 + 1469;
    int b = C56 * 7;
    int c = (1469 + C61) % 10000;
    boolean flag = false;
    string label = "rec1469";
};
type R1470 record {
    int id;
    int a = C62 + 1470;
    int b = C59 * 1;
    int c = (1470 + C62) % 10000;
    boolean flag = true;
    string label = "rec1470";
};
type R1471 record {
    int id;
    int a = C63 + 1471;
    int b = C62 * 2;
    int c = (1471 + C63) % 10000;
    boolean flag = false;
    string label = "rec1471";
};
type R1472 record {
    int id;
    int a = C0 + 1472;
    int b = C1 * 3;
    int c = (1472 + C0) % 10000;
    boolean flag = true;
    string label = "rec1472";
};
type R1473 record {
    int id;
    int a = C1 + 1473;
    int b = C4 * 4;
    int c = (1473 + C1) % 10000;
    boolean flag = false;
    string label = "rec1473";
};
type R1474 record {
    int id;
    int a = C2 + 1474;
    int b = C7 * 5;
    int c = (1474 + C2) % 10000;
    boolean flag = true;
    string label = "rec1474";
};
type R1475 record {
    int id;
    int a = C3 + 1475;
    int b = C10 * 6;
    int c = (1475 + C3) % 10000;
    boolean flag = false;
    string label = "rec1475";
};
type R1476 record {
    int id;
    int a = C4 + 1476;
    int b = C13 * 7;
    int c = (1476 + C4) % 10000;
    boolean flag = true;
    string label = "rec1476";
};
type R1477 record {
    int id;
    int a = C5 + 1477;
    int b = C16 * 1;
    int c = (1477 + C5) % 10000;
    boolean flag = false;
    string label = "rec1477";
};
type R1478 record {
    int id;
    int a = C6 + 1478;
    int b = C19 * 2;
    int c = (1478 + C6) % 10000;
    boolean flag = true;
    string label = "rec1478";
};
type R1479 record {
    int id;
    int a = C7 + 1479;
    int b = C22 * 3;
    int c = (1479 + C7) % 10000;
    boolean flag = false;
    string label = "rec1479";
};
type R1480 record {
    int id;
    int a = C8 + 1480;
    int b = C25 * 4;
    int c = (1480 + C8) % 10000;
    boolean flag = true;
    string label = "rec1480";
};
type R1481 record {
    int id;
    int a = C9 + 1481;
    int b = C28 * 5;
    int c = (1481 + C9) % 10000;
    boolean flag = false;
    string label = "rec1481";
};
type R1482 record {
    int id;
    int a = C10 + 1482;
    int b = C31 * 6;
    int c = (1482 + C10) % 10000;
    boolean flag = true;
    string label = "rec1482";
};
type R1483 record {
    int id;
    int a = C11 + 1483;
    int b = C34 * 7;
    int c = (1483 + C11) % 10000;
    boolean flag = false;
    string label = "rec1483";
};
type R1484 record {
    int id;
    int a = C12 + 1484;
    int b = C37 * 1;
    int c = (1484 + C12) % 10000;
    boolean flag = true;
    string label = "rec1484";
};
type R1485 record {
    int id;
    int a = C13 + 1485;
    int b = C40 * 2;
    int c = (1485 + C13) % 10000;
    boolean flag = false;
    string label = "rec1485";
};
type R1486 record {
    int id;
    int a = C14 + 1486;
    int b = C43 * 3;
    int c = (1486 + C14) % 10000;
    boolean flag = true;
    string label = "rec1486";
};
type R1487 record {
    int id;
    int a = C15 + 1487;
    int b = C46 * 4;
    int c = (1487 + C15) % 10000;
    boolean flag = false;
    string label = "rec1487";
};
type R1488 record {
    int id;
    int a = C16 + 1488;
    int b = C49 * 5;
    int c = (1488 + C16) % 10000;
    boolean flag = true;
    string label = "rec1488";
};
type R1489 record {
    int id;
    int a = C17 + 1489;
    int b = C52 * 6;
    int c = (1489 + C17) % 10000;
    boolean flag = false;
    string label = "rec1489";
};
type R1490 record {
    int id;
    int a = C18 + 1490;
    int b = C55 * 7;
    int c = (1490 + C18) % 10000;
    boolean flag = true;
    string label = "rec1490";
};
type R1491 record {
    int id;
    int a = C19 + 1491;
    int b = C58 * 1;
    int c = (1491 + C19) % 10000;
    boolean flag = false;
    string label = "rec1491";
};
type R1492 record {
    int id;
    int a = C20 + 1492;
    int b = C61 * 2;
    int c = (1492 + C20) % 10000;
    boolean flag = true;
    string label = "rec1492";
};
type R1493 record {
    int id;
    int a = C21 + 1493;
    int b = C0 * 3;
    int c = (1493 + C21) % 10000;
    boolean flag = false;
    string label = "rec1493";
};
type R1494 record {
    int id;
    int a = C22 + 1494;
    int b = C3 * 4;
    int c = (1494 + C22) % 10000;
    boolean flag = true;
    string label = "rec1494";
};
type R1495 record {
    int id;
    int a = C23 + 1495;
    int b = C6 * 5;
    int c = (1495 + C23) % 10000;
    boolean flag = false;
    string label = "rec1495";
};
type R1496 record {
    int id;
    int a = C24 + 1496;
    int b = C9 * 6;
    int c = (1496 + C24) % 10000;
    boolean flag = true;
    string label = "rec1496";
};
type R1497 record {
    int id;
    int a = C25 + 1497;
    int b = C12 * 7;
    int c = (1497 + C25) % 10000;
    boolean flag = false;
    string label = "rec1497";
};
type R1498 record {
    int id;
    int a = C26 + 1498;
    int b = C15 * 1;
    int c = (1498 + C26) % 10000;
    boolean flag = true;
    string label = "rec1498";
};
type R1499 record {
    int id;
    int a = C27 + 1499;
    int b = C18 * 2;
    int c = (1499 + C27) % 10000;
    boolean flag = false;
    string label = "rec1499";
};
type R1500 record {
    int id;
    int a = C28 + 1500;
    int b = C21 * 3;
    int c = (1500 + C28) % 10000;
    boolean flag = true;
    string label = "rec1500";
};
type R1501 record {
    int id;
    int a = C29 + 1501;
    int b = C24 * 4;
    int c = (1501 + C29) % 10000;
    boolean flag = false;
    string label = "rec1501";
};
type R1502 record {
    int id;
    int a = C30 + 1502;
    int b = C27 * 5;
    int c = (1502 + C30) % 10000;
    boolean flag = true;
    string label = "rec1502";
};
type R1503 record {
    int id;
    int a = C31 + 1503;
    int b = C30 * 6;
    int c = (1503 + C31) % 10000;
    boolean flag = false;
    string label = "rec1503";
};
type R1504 record {
    int id;
    int a = C32 + 1504;
    int b = C33 * 7;
    int c = (1504 + C32) % 10000;
    boolean flag = true;
    string label = "rec1504";
};
type R1505 record {
    int id;
    int a = C33 + 1505;
    int b = C36 * 1;
    int c = (1505 + C33) % 10000;
    boolean flag = false;
    string label = "rec1505";
};
type R1506 record {
    int id;
    int a = C34 + 1506;
    int b = C39 * 2;
    int c = (1506 + C34) % 10000;
    boolean flag = true;
    string label = "rec1506";
};
type R1507 record {
    int id;
    int a = C35 + 1507;
    int b = C42 * 3;
    int c = (1507 + C35) % 10000;
    boolean flag = false;
    string label = "rec1507";
};
type R1508 record {
    int id;
    int a = C36 + 1508;
    int b = C45 * 4;
    int c = (1508 + C36) % 10000;
    boolean flag = true;
    string label = "rec1508";
};
type R1509 record {
    int id;
    int a = C37 + 1509;
    int b = C48 * 5;
    int c = (1509 + C37) % 10000;
    boolean flag = false;
    string label = "rec1509";
};
type R1510 record {
    int id;
    int a = C38 + 1510;
    int b = C51 * 6;
    int c = (1510 + C38) % 10000;
    boolean flag = true;
    string label = "rec1510";
};
type R1511 record {
    int id;
    int a = C39 + 1511;
    int b = C54 * 7;
    int c = (1511 + C39) % 10000;
    boolean flag = false;
    string label = "rec1511";
};
type R1512 record {
    int id;
    int a = C40 + 1512;
    int b = C57 * 1;
    int c = (1512 + C40) % 10000;
    boolean flag = true;
    string label = "rec1512";
};
type R1513 record {
    int id;
    int a = C41 + 1513;
    int b = C60 * 2;
    int c = (1513 + C41) % 10000;
    boolean flag = false;
    string label = "rec1513";
};
type R1514 record {
    int id;
    int a = C42 + 1514;
    int b = C63 * 3;
    int c = (1514 + C42) % 10000;
    boolean flag = true;
    string label = "rec1514";
};
type R1515 record {
    int id;
    int a = C43 + 1515;
    int b = C2 * 4;
    int c = (1515 + C43) % 10000;
    boolean flag = false;
    string label = "rec1515";
};
type R1516 record {
    int id;
    int a = C44 + 1516;
    int b = C5 * 5;
    int c = (1516 + C44) % 10000;
    boolean flag = true;
    string label = "rec1516";
};
type R1517 record {
    int id;
    int a = C45 + 1517;
    int b = C8 * 6;
    int c = (1517 + C45) % 10000;
    boolean flag = false;
    string label = "rec1517";
};
type R1518 record {
    int id;
    int a = C46 + 1518;
    int b = C11 * 7;
    int c = (1518 + C46) % 10000;
    boolean flag = true;
    string label = "rec1518";
};
type R1519 record {
    int id;
    int a = C47 + 1519;
    int b = C14 * 1;
    int c = (1519 + C47) % 10000;
    boolean flag = false;
    string label = "rec1519";
};
type R1520 record {
    int id;
    int a = C48 + 1520;
    int b = C17 * 2;
    int c = (1520 + C48) % 10000;
    boolean flag = true;
    string label = "rec1520";
};
type R1521 record {
    int id;
    int a = C49 + 1521;
    int b = C20 * 3;
    int c = (1521 + C49) % 10000;
    boolean flag = false;
    string label = "rec1521";
};
type R1522 record {
    int id;
    int a = C50 + 1522;
    int b = C23 * 4;
    int c = (1522 + C50) % 10000;
    boolean flag = true;
    string label = "rec1522";
};
type R1523 record {
    int id;
    int a = C51 + 1523;
    int b = C26 * 5;
    int c = (1523 + C51) % 10000;
    boolean flag = false;
    string label = "rec1523";
};
type R1524 record {
    int id;
    int a = C52 + 1524;
    int b = C29 * 6;
    int c = (1524 + C52) % 10000;
    boolean flag = true;
    string label = "rec1524";
};
type R1525 record {
    int id;
    int a = C53 + 1525;
    int b = C32 * 7;
    int c = (1525 + C53) % 10000;
    boolean flag = false;
    string label = "rec1525";
};
type R1526 record {
    int id;
    int a = C54 + 1526;
    int b = C35 * 1;
    int c = (1526 + C54) % 10000;
    boolean flag = true;
    string label = "rec1526";
};
type R1527 record {
    int id;
    int a = C55 + 1527;
    int b = C38 * 2;
    int c = (1527 + C55) % 10000;
    boolean flag = false;
    string label = "rec1527";
};
type R1528 record {
    int id;
    int a = C56 + 1528;
    int b = C41 * 3;
    int c = (1528 + C56) % 10000;
    boolean flag = true;
    string label = "rec1528";
};
type R1529 record {
    int id;
    int a = C57 + 1529;
    int b = C44 * 4;
    int c = (1529 + C57) % 10000;
    boolean flag = false;
    string label = "rec1529";
};
type R1530 record {
    int id;
    int a = C58 + 1530;
    int b = C47 * 5;
    int c = (1530 + C58) % 10000;
    boolean flag = true;
    string label = "rec1530";
};
type R1531 record {
    int id;
    int a = C59 + 1531;
    int b = C50 * 6;
    int c = (1531 + C59) % 10000;
    boolean flag = false;
    string label = "rec1531";
};
type R1532 record {
    int id;
    int a = C60 + 1532;
    int b = C53 * 7;
    int c = (1532 + C60) % 10000;
    boolean flag = true;
    string label = "rec1532";
};
type R1533 record {
    int id;
    int a = C61 + 1533;
    int b = C56 * 1;
    int c = (1533 + C61) % 10000;
    boolean flag = false;
    string label = "rec1533";
};
type R1534 record {
    int id;
    int a = C62 + 1534;
    int b = C59 * 2;
    int c = (1534 + C62) % 10000;
    boolean flag = true;
    string label = "rec1534";
};
type R1535 record {
    int id;
    int a = C63 + 1535;
    int b = C62 * 3;
    int c = (1535 + C63) % 10000;
    boolean flag = false;
    string label = "rec1535";
};
type R1536 record {
    int id;
    int a = C0 + 1536;
    int b = C1 * 4;
    int c = (1536 + C0) % 10000;
    boolean flag = true;
    string label = "rec1536";
};
type R1537 record {
    int id;
    int a = C1 + 1537;
    int b = C4 * 5;
    int c = (1537 + C1) % 10000;
    boolean flag = false;
    string label = "rec1537";
};
type R1538 record {
    int id;
    int a = C2 + 1538;
    int b = C7 * 6;
    int c = (1538 + C2) % 10000;
    boolean flag = true;
    string label = "rec1538";
};
type R1539 record {
    int id;
    int a = C3 + 1539;
    int b = C10 * 7;
    int c = (1539 + C3) % 10000;
    boolean flag = false;
    string label = "rec1539";
};
type R1540 record {
    int id;
    int a = C4 + 1540;
    int b = C13 * 1;
    int c = (1540 + C4) % 10000;
    boolean flag = true;
    string label = "rec1540";
};
type R1541 record {
    int id;
    int a = C5 + 1541;
    int b = C16 * 2;
    int c = (1541 + C5) % 10000;
    boolean flag = false;
    string label = "rec1541";
};
type R1542 record {
    int id;
    int a = C6 + 1542;
    int b = C19 * 3;
    int c = (1542 + C6) % 10000;
    boolean flag = true;
    string label = "rec1542";
};
type R1543 record {
    int id;
    int a = C7 + 1543;
    int b = C22 * 4;
    int c = (1543 + C7) % 10000;
    boolean flag = false;
    string label = "rec1543";
};
type R1544 record {
    int id;
    int a = C8 + 1544;
    int b = C25 * 5;
    int c = (1544 + C8) % 10000;
    boolean flag = true;
    string label = "rec1544";
};
type R1545 record {
    int id;
    int a = C9 + 1545;
    int b = C28 * 6;
    int c = (1545 + C9) % 10000;
    boolean flag = false;
    string label = "rec1545";
};
type R1546 record {
    int id;
    int a = C10 + 1546;
    int b = C31 * 7;
    int c = (1546 + C10) % 10000;
    boolean flag = true;
    string label = "rec1546";
};
type R1547 record {
    int id;
    int a = C11 + 1547;
    int b = C34 * 1;
    int c = (1547 + C11) % 10000;
    boolean flag = false;
    string label = "rec1547";
};
type R1548 record {
    int id;
    int a = C12 + 1548;
    int b = C37 * 2;
    int c = (1548 + C12) % 10000;
    boolean flag = true;
    string label = "rec1548";
};
type R1549 record {
    int id;
    int a = C13 + 1549;
    int b = C40 * 3;
    int c = (1549 + C13) % 10000;
    boolean flag = false;
    string label = "rec1549";
};
type R1550 record {
    int id;
    int a = C14 + 1550;
    int b = C43 * 4;
    int c = (1550 + C14) % 10000;
    boolean flag = true;
    string label = "rec1550";
};
type R1551 record {
    int id;
    int a = C15 + 1551;
    int b = C46 * 5;
    int c = (1551 + C15) % 10000;
    boolean flag = false;
    string label = "rec1551";
};
type R1552 record {
    int id;
    int a = C16 + 1552;
    int b = C49 * 6;
    int c = (1552 + C16) % 10000;
    boolean flag = true;
    string label = "rec1552";
};
type R1553 record {
    int id;
    int a = C17 + 1553;
    int b = C52 * 7;
    int c = (1553 + C17) % 10000;
    boolean flag = false;
    string label = "rec1553";
};
type R1554 record {
    int id;
    int a = C18 + 1554;
    int b = C55 * 1;
    int c = (1554 + C18) % 10000;
    boolean flag = true;
    string label = "rec1554";
};
type R1555 record {
    int id;
    int a = C19 + 1555;
    int b = C58 * 2;
    int c = (1555 + C19) % 10000;
    boolean flag = false;
    string label = "rec1555";
};
type R1556 record {
    int id;
    int a = C20 + 1556;
    int b = C61 * 3;
    int c = (1556 + C20) % 10000;
    boolean flag = true;
    string label = "rec1556";
};
type R1557 record {
    int id;
    int a = C21 + 1557;
    int b = C0 * 4;
    int c = (1557 + C21) % 10000;
    boolean flag = false;
    string label = "rec1557";
};
type R1558 record {
    int id;
    int a = C22 + 1558;
    int b = C3 * 5;
    int c = (1558 + C22) % 10000;
    boolean flag = true;
    string label = "rec1558";
};
type R1559 record {
    int id;
    int a = C23 + 1559;
    int b = C6 * 6;
    int c = (1559 + C23) % 10000;
    boolean flag = false;
    string label = "rec1559";
};
type R1560 record {
    int id;
    int a = C24 + 1560;
    int b = C9 * 7;
    int c = (1560 + C24) % 10000;
    boolean flag = true;
    string label = "rec1560";
};
type R1561 record {
    int id;
    int a = C25 + 1561;
    int b = C12 * 1;
    int c = (1561 + C25) % 10000;
    boolean flag = false;
    string label = "rec1561";
};
type R1562 record {
    int id;
    int a = C26 + 1562;
    int b = C15 * 2;
    int c = (1562 + C26) % 10000;
    boolean flag = true;
    string label = "rec1562";
};
type R1563 record {
    int id;
    int a = C27 + 1563;
    int b = C18 * 3;
    int c = (1563 + C27) % 10000;
    boolean flag = false;
    string label = "rec1563";
};
type R1564 record {
    int id;
    int a = C28 + 1564;
    int b = C21 * 4;
    int c = (1564 + C28) % 10000;
    boolean flag = true;
    string label = "rec1564";
};
type R1565 record {
    int id;
    int a = C29 + 1565;
    int b = C24 * 5;
    int c = (1565 + C29) % 10000;
    boolean flag = false;
    string label = "rec1565";
};
type R1566 record {
    int id;
    int a = C30 + 1566;
    int b = C27 * 6;
    int c = (1566 + C30) % 10000;
    boolean flag = true;
    string label = "rec1566";
};
type R1567 record {
    int id;
    int a = C31 + 1567;
    int b = C30 * 7;
    int c = (1567 + C31) % 10000;
    boolean flag = false;
    string label = "rec1567";
};
type R1568 record {
    int id;
    int a = C32 + 1568;
    int b = C33 * 1;
    int c = (1568 + C32) % 10000;
    boolean flag = true;
    string label = "rec1568";
};
type R1569 record {
    int id;
    int a = C33 + 1569;
    int b = C36 * 2;
    int c = (1569 + C33) % 10000;
    boolean flag = false;
    string label = "rec1569";
};
type R1570 record {
    int id;
    int a = C34 + 1570;
    int b = C39 * 3;
    int c = (1570 + C34) % 10000;
    boolean flag = true;
    string label = "rec1570";
};
type R1571 record {
    int id;
    int a = C35 + 1571;
    int b = C42 * 4;
    int c = (1571 + C35) % 10000;
    boolean flag = false;
    string label = "rec1571";
};
type R1572 record {
    int id;
    int a = C36 + 1572;
    int b = C45 * 5;
    int c = (1572 + C36) % 10000;
    boolean flag = true;
    string label = "rec1572";
};
type R1573 record {
    int id;
    int a = C37 + 1573;
    int b = C48 * 6;
    int c = (1573 + C37) % 10000;
    boolean flag = false;
    string label = "rec1573";
};
type R1574 record {
    int id;
    int a = C38 + 1574;
    int b = C51 * 7;
    int c = (1574 + C38) % 10000;
    boolean flag = true;
    string label = "rec1574";
};
type R1575 record {
    int id;
    int a = C39 + 1575;
    int b = C54 * 1;
    int c = (1575 + C39) % 10000;
    boolean flag = false;
    string label = "rec1575";
};
type R1576 record {
    int id;
    int a = C40 + 1576;
    int b = C57 * 2;
    int c = (1576 + C40) % 10000;
    boolean flag = true;
    string label = "rec1576";
};
type R1577 record {
    int id;
    int a = C41 + 1577;
    int b = C60 * 3;
    int c = (1577 + C41) % 10000;
    boolean flag = false;
    string label = "rec1577";
};
type R1578 record {
    int id;
    int a = C42 + 1578;
    int b = C63 * 4;
    int c = (1578 + C42) % 10000;
    boolean flag = true;
    string label = "rec1578";
};
type R1579 record {
    int id;
    int a = C43 + 1579;
    int b = C2 * 5;
    int c = (1579 + C43) % 10000;
    boolean flag = false;
    string label = "rec1579";
};
type R1580 record {
    int id;
    int a = C44 + 1580;
    int b = C5 * 6;
    int c = (1580 + C44) % 10000;
    boolean flag = true;
    string label = "rec1580";
};
type R1581 record {
    int id;
    int a = C45 + 1581;
    int b = C8 * 7;
    int c = (1581 + C45) % 10000;
    boolean flag = false;
    string label = "rec1581";
};
type R1582 record {
    int id;
    int a = C46 + 1582;
    int b = C11 * 1;
    int c = (1582 + C46) % 10000;
    boolean flag = true;
    string label = "rec1582";
};
type R1583 record {
    int id;
    int a = C47 + 1583;
    int b = C14 * 2;
    int c = (1583 + C47) % 10000;
    boolean flag = false;
    string label = "rec1583";
};
type R1584 record {
    int id;
    int a = C48 + 1584;
    int b = C17 * 3;
    int c = (1584 + C48) % 10000;
    boolean flag = true;
    string label = "rec1584";
};
type R1585 record {
    int id;
    int a = C49 + 1585;
    int b = C20 * 4;
    int c = (1585 + C49) % 10000;
    boolean flag = false;
    string label = "rec1585";
};
type R1586 record {
    int id;
    int a = C50 + 1586;
    int b = C23 * 5;
    int c = (1586 + C50) % 10000;
    boolean flag = true;
    string label = "rec1586";
};
type R1587 record {
    int id;
    int a = C51 + 1587;
    int b = C26 * 6;
    int c = (1587 + C51) % 10000;
    boolean flag = false;
    string label = "rec1587";
};
type R1588 record {
    int id;
    int a = C52 + 1588;
    int b = C29 * 7;
    int c = (1588 + C52) % 10000;
    boolean flag = true;
    string label = "rec1588";
};
type R1589 record {
    int id;
    int a = C53 + 1589;
    int b = C32 * 1;
    int c = (1589 + C53) % 10000;
    boolean flag = false;
    string label = "rec1589";
};
type R1590 record {
    int id;
    int a = C54 + 1590;
    int b = C35 * 2;
    int c = (1590 + C54) % 10000;
    boolean flag = true;
    string label = "rec1590";
};
type R1591 record {
    int id;
    int a = C55 + 1591;
    int b = C38 * 3;
    int c = (1591 + C55) % 10000;
    boolean flag = false;
    string label = "rec1591";
};
type R1592 record {
    int id;
    int a = C56 + 1592;
    int b = C41 * 4;
    int c = (1592 + C56) % 10000;
    boolean flag = true;
    string label = "rec1592";
};
type R1593 record {
    int id;
    int a = C57 + 1593;
    int b = C44 * 5;
    int c = (1593 + C57) % 10000;
    boolean flag = false;
    string label = "rec1593";
};
type R1594 record {
    int id;
    int a = C58 + 1594;
    int b = C47 * 6;
    int c = (1594 + C58) % 10000;
    boolean flag = true;
    string label = "rec1594";
};
type R1595 record {
    int id;
    int a = C59 + 1595;
    int b = C50 * 7;
    int c = (1595 + C59) % 10000;
    boolean flag = false;
    string label = "rec1595";
};
type R1596 record {
    int id;
    int a = C60 + 1596;
    int b = C53 * 1;
    int c = (1596 + C60) % 10000;
    boolean flag = true;
    string label = "rec1596";
};
type R1597 record {
    int id;
    int a = C61 + 1597;
    int b = C56 * 2;
    int c = (1597 + C61) % 10000;
    boolean flag = false;
    string label = "rec1597";
};
type R1598 record {
    int id;
    int a = C62 + 1598;
    int b = C59 * 3;
    int c = (1598 + C62) % 10000;
    boolean flag = true;
    string label = "rec1598";
};
type R1599 record {
    int id;
    int a = C63 + 1599;
    int b = C62 * 4;
    int c = (1599 + C63) % 10000;
    boolean flag = false;
    string label = "rec1599";
};
type R1600 record {
    int id;
    int a = C0 + 1600;
    int b = C1 * 5;
    int c = (1600 + C0) % 10000;
    boolean flag = true;
    string label = "rec1600";
};
type R1601 record {
    int id;
    int a = C1 + 1601;
    int b = C4 * 6;
    int c = (1601 + C1) % 10000;
    boolean flag = false;
    string label = "rec1601";
};
type R1602 record {
    int id;
    int a = C2 + 1602;
    int b = C7 * 7;
    int c = (1602 + C2) % 10000;
    boolean flag = true;
    string label = "rec1602";
};
type R1603 record {
    int id;
    int a = C3 + 1603;
    int b = C10 * 1;
    int c = (1603 + C3) % 10000;
    boolean flag = false;
    string label = "rec1603";
};
type R1604 record {
    int id;
    int a = C4 + 1604;
    int b = C13 * 2;
    int c = (1604 + C4) % 10000;
    boolean flag = true;
    string label = "rec1604";
};
type R1605 record {
    int id;
    int a = C5 + 1605;
    int b = C16 * 3;
    int c = (1605 + C5) % 10000;
    boolean flag = false;
    string label = "rec1605";
};
type R1606 record {
    int id;
    int a = C6 + 1606;
    int b = C19 * 4;
    int c = (1606 + C6) % 10000;
    boolean flag = true;
    string label = "rec1606";
};
type R1607 record {
    int id;
    int a = C7 + 1607;
    int b = C22 * 5;
    int c = (1607 + C7) % 10000;
    boolean flag = false;
    string label = "rec1607";
};
type R1608 record {
    int id;
    int a = C8 + 1608;
    int b = C25 * 6;
    int c = (1608 + C8) % 10000;
    boolean flag = true;
    string label = "rec1608";
};
type R1609 record {
    int id;
    int a = C9 + 1609;
    int b = C28 * 7;
    int c = (1609 + C9) % 10000;
    boolean flag = false;
    string label = "rec1609";
};
type R1610 record {
    int id;
    int a = C10 + 1610;
    int b = C31 * 1;
    int c = (1610 + C10) % 10000;
    boolean flag = true;
    string label = "rec1610";
};
type R1611 record {
    int id;
    int a = C11 + 1611;
    int b = C34 * 2;
    int c = (1611 + C11) % 10000;
    boolean flag = false;
    string label = "rec1611";
};
type R1612 record {
    int id;
    int a = C12 + 1612;
    int b = C37 * 3;
    int c = (1612 + C12) % 10000;
    boolean flag = true;
    string label = "rec1612";
};
type R1613 record {
    int id;
    int a = C13 + 1613;
    int b = C40 * 4;
    int c = (1613 + C13) % 10000;
    boolean flag = false;
    string label = "rec1613";
};
type R1614 record {
    int id;
    int a = C14 + 1614;
    int b = C43 * 5;
    int c = (1614 + C14) % 10000;
    boolean flag = true;
    string label = "rec1614";
};
type R1615 record {
    int id;
    int a = C15 + 1615;
    int b = C46 * 6;
    int c = (1615 + C15) % 10000;
    boolean flag = false;
    string label = "rec1615";
};
type R1616 record {
    int id;
    int a = C16 + 1616;
    int b = C49 * 7;
    int c = (1616 + C16) % 10000;
    boolean flag = true;
    string label = "rec1616";
};
type R1617 record {
    int id;
    int a = C17 + 1617;
    int b = C52 * 1;
    int c = (1617 + C17) % 10000;
    boolean flag = false;
    string label = "rec1617";
};
type R1618 record {
    int id;
    int a = C18 + 1618;
    int b = C55 * 2;
    int c = (1618 + C18) % 10000;
    boolean flag = true;
    string label = "rec1618";
};
type R1619 record {
    int id;
    int a = C19 + 1619;
    int b = C58 * 3;
    int c = (1619 + C19) % 10000;
    boolean flag = false;
    string label = "rec1619";
};
type R1620 record {
    int id;
    int a = C20 + 1620;
    int b = C61 * 4;
    int c = (1620 + C20) % 10000;
    boolean flag = true;
    string label = "rec1620";
};
type R1621 record {
    int id;
    int a = C21 + 1621;
    int b = C0 * 5;
    int c = (1621 + C21) % 10000;
    boolean flag = false;
    string label = "rec1621";
};
type R1622 record {
    int id;
    int a = C22 + 1622;
    int b = C3 * 6;
    int c = (1622 + C22) % 10000;
    boolean flag = true;
    string label = "rec1622";
};
type R1623 record {
    int id;
    int a = C23 + 1623;
    int b = C6 * 7;
    int c = (1623 + C23) % 10000;
    boolean flag = false;
    string label = "rec1623";
};
type R1624 record {
    int id;
    int a = C24 + 1624;
    int b = C9 * 1;
    int c = (1624 + C24) % 10000;
    boolean flag = true;
    string label = "rec1624";
};
type R1625 record {
    int id;
    int a = C25 + 1625;
    int b = C12 * 2;
    int c = (1625 + C25) % 10000;
    boolean flag = false;
    string label = "rec1625";
};
type R1626 record {
    int id;
    int a = C26 + 1626;
    int b = C15 * 3;
    int c = (1626 + C26) % 10000;
    boolean flag = true;
    string label = "rec1626";
};
type R1627 record {
    int id;
    int a = C27 + 1627;
    int b = C18 * 4;
    int c = (1627 + C27) % 10000;
    boolean flag = false;
    string label = "rec1627";
};
type R1628 record {
    int id;
    int a = C28 + 1628;
    int b = C21 * 5;
    int c = (1628 + C28) % 10000;
    boolean flag = true;
    string label = "rec1628";
};
type R1629 record {
    int id;
    int a = C29 + 1629;
    int b = C24 * 6;
    int c = (1629 + C29) % 10000;
    boolean flag = false;
    string label = "rec1629";
};
type R1630 record {
    int id;
    int a = C30 + 1630;
    int b = C27 * 7;
    int c = (1630 + C30) % 10000;
    boolean flag = true;
    string label = "rec1630";
};
type R1631 record {
    int id;
    int a = C31 + 1631;
    int b = C30 * 1;
    int c = (1631 + C31) % 10000;
    boolean flag = false;
    string label = "rec1631";
};
type R1632 record {
    int id;
    int a = C32 + 1632;
    int b = C33 * 2;
    int c = (1632 + C32) % 10000;
    boolean flag = true;
    string label = "rec1632";
};
type R1633 record {
    int id;
    int a = C33 + 1633;
    int b = C36 * 3;
    int c = (1633 + C33) % 10000;
    boolean flag = false;
    string label = "rec1633";
};
type R1634 record {
    int id;
    int a = C34 + 1634;
    int b = C39 * 4;
    int c = (1634 + C34) % 10000;
    boolean flag = true;
    string label = "rec1634";
};
type R1635 record {
    int id;
    int a = C35 + 1635;
    int b = C42 * 5;
    int c = (1635 + C35) % 10000;
    boolean flag = false;
    string label = "rec1635";
};
type R1636 record {
    int id;
    int a = C36 + 1636;
    int b = C45 * 6;
    int c = (1636 + C36) % 10000;
    boolean flag = true;
    string label = "rec1636";
};
type R1637 record {
    int id;
    int a = C37 + 1637;
    int b = C48 * 7;
    int c = (1637 + C37) % 10000;
    boolean flag = false;
    string label = "rec1637";
};
type R1638 record {
    int id;
    int a = C38 + 1638;
    int b = C51 * 1;
    int c = (1638 + C38) % 10000;
    boolean flag = true;
    string label = "rec1638";
};
type R1639 record {
    int id;
    int a = C39 + 1639;
    int b = C54 * 2;
    int c = (1639 + C39) % 10000;
    boolean flag = false;
    string label = "rec1639";
};
type R1640 record {
    int id;
    int a = C40 + 1640;
    int b = C57 * 3;
    int c = (1640 + C40) % 10000;
    boolean flag = true;
    string label = "rec1640";
};
type R1641 record {
    int id;
    int a = C41 + 1641;
    int b = C60 * 4;
    int c = (1641 + C41) % 10000;
    boolean flag = false;
    string label = "rec1641";
};
type R1642 record {
    int id;
    int a = C42 + 1642;
    int b = C63 * 5;
    int c = (1642 + C42) % 10000;
    boolean flag = true;
    string label = "rec1642";
};
type R1643 record {
    int id;
    int a = C43 + 1643;
    int b = C2 * 6;
    int c = (1643 + C43) % 10000;
    boolean flag = false;
    string label = "rec1643";
};
type R1644 record {
    int id;
    int a = C44 + 1644;
    int b = C5 * 7;
    int c = (1644 + C44) % 10000;
    boolean flag = true;
    string label = "rec1644";
};
type R1645 record {
    int id;
    int a = C45 + 1645;
    int b = C8 * 1;
    int c = (1645 + C45) % 10000;
    boolean flag = false;
    string label = "rec1645";
};
type R1646 record {
    int id;
    int a = C46 + 1646;
    int b = C11 * 2;
    int c = (1646 + C46) % 10000;
    boolean flag = true;
    string label = "rec1646";
};
type R1647 record {
    int id;
    int a = C47 + 1647;
    int b = C14 * 3;
    int c = (1647 + C47) % 10000;
    boolean flag = false;
    string label = "rec1647";
};
type R1648 record {
    int id;
    int a = C48 + 1648;
    int b = C17 * 4;
    int c = (1648 + C48) % 10000;
    boolean flag = true;
    string label = "rec1648";
};
type R1649 record {
    int id;
    int a = C49 + 1649;
    int b = C20 * 5;
    int c = (1649 + C49) % 10000;
    boolean flag = false;
    string label = "rec1649";
};
type R1650 record {
    int id;
    int a = C50 + 1650;
    int b = C23 * 6;
    int c = (1650 + C50) % 10000;
    boolean flag = true;
    string label = "rec1650";
};
type R1651 record {
    int id;
    int a = C51 + 1651;
    int b = C26 * 7;
    int c = (1651 + C51) % 10000;
    boolean flag = false;
    string label = "rec1651";
};
type R1652 record {
    int id;
    int a = C52 + 1652;
    int b = C29 * 1;
    int c = (1652 + C52) % 10000;
    boolean flag = true;
    string label = "rec1652";
};
type R1653 record {
    int id;
    int a = C53 + 1653;
    int b = C32 * 2;
    int c = (1653 + C53) % 10000;
    boolean flag = false;
    string label = "rec1653";
};
type R1654 record {
    int id;
    int a = C54 + 1654;
    int b = C35 * 3;
    int c = (1654 + C54) % 10000;
    boolean flag = true;
    string label = "rec1654";
};
type R1655 record {
    int id;
    int a = C55 + 1655;
    int b = C38 * 4;
    int c = (1655 + C55) % 10000;
    boolean flag = false;
    string label = "rec1655";
};
type R1656 record {
    int id;
    int a = C56 + 1656;
    int b = C41 * 5;
    int c = (1656 + C56) % 10000;
    boolean flag = true;
    string label = "rec1656";
};
type R1657 record {
    int id;
    int a = C57 + 1657;
    int b = C44 * 6;
    int c = (1657 + C57) % 10000;
    boolean flag = false;
    string label = "rec1657";
};
type R1658 record {
    int id;
    int a = C58 + 1658;
    int b = C47 * 7;
    int c = (1658 + C58) % 10000;
    boolean flag = true;
    string label = "rec1658";
};
type R1659 record {
    int id;
    int a = C59 + 1659;
    int b = C50 * 1;
    int c = (1659 + C59) % 10000;
    boolean flag = false;
    string label = "rec1659";
};
type R1660 record {
    int id;
    int a = C60 + 1660;
    int b = C53 * 2;
    int c = (1660 + C60) % 10000;
    boolean flag = true;
    string label = "rec1660";
};
type R1661 record {
    int id;
    int a = C61 + 1661;
    int b = C56 * 3;
    int c = (1661 + C61) % 10000;
    boolean flag = false;
    string label = "rec1661";
};
type R1662 record {
    int id;
    int a = C62 + 1662;
    int b = C59 * 4;
    int c = (1662 + C62) % 10000;
    boolean flag = true;
    string label = "rec1662";
};
type R1663 record {
    int id;
    int a = C63 + 1663;
    int b = C62 * 5;
    int c = (1663 + C63) % 10000;
    boolean flag = false;
    string label = "rec1663";
};
type R1664 record {
    int id;
    int a = C0 + 1664;
    int b = C1 * 6;
    int c = (1664 + C0) % 10000;
    boolean flag = true;
    string label = "rec1664";
};
type R1665 record {
    int id;
    int a = C1 + 1665;
    int b = C4 * 7;
    int c = (1665 + C1) % 10000;
    boolean flag = false;
    string label = "rec1665";
};
type R1666 record {
    int id;
    int a = C2 + 1666;
    int b = C7 * 1;
    int c = (1666 + C2) % 10000;
    boolean flag = true;
    string label = "rec1666";
};
type R1667 record {
    int id;
    int a = C3 + 1667;
    int b = C10 * 2;
    int c = (1667 + C3) % 10000;
    boolean flag = false;
    string label = "rec1667";
};
type R1668 record {
    int id;
    int a = C4 + 1668;
    int b = C13 * 3;
    int c = (1668 + C4) % 10000;
    boolean flag = true;
    string label = "rec1668";
};
type R1669 record {
    int id;
    int a = C5 + 1669;
    int b = C16 * 4;
    int c = (1669 + C5) % 10000;
    boolean flag = false;
    string label = "rec1669";
};
type R1670 record {
    int id;
    int a = C6 + 1670;
    int b = C19 * 5;
    int c = (1670 + C6) % 10000;
    boolean flag = true;
    string label = "rec1670";
};
type R1671 record {
    int id;
    int a = C7 + 1671;
    int b = C22 * 6;
    int c = (1671 + C7) % 10000;
    boolean flag = false;
    string label = "rec1671";
};
type R1672 record {
    int id;
    int a = C8 + 1672;
    int b = C25 * 7;
    int c = (1672 + C8) % 10000;
    boolean flag = true;
    string label = "rec1672";
};
type R1673 record {
    int id;
    int a = C9 + 1673;
    int b = C28 * 1;
    int c = (1673 + C9) % 10000;
    boolean flag = false;
    string label = "rec1673";
};
type R1674 record {
    int id;
    int a = C10 + 1674;
    int b = C31 * 2;
    int c = (1674 + C10) % 10000;
    boolean flag = true;
    string label = "rec1674";
};
type R1675 record {
    int id;
    int a = C11 + 1675;
    int b = C34 * 3;
    int c = (1675 + C11) % 10000;
    boolean flag = false;
    string label = "rec1675";
};
type R1676 record {
    int id;
    int a = C12 + 1676;
    int b = C37 * 4;
    int c = (1676 + C12) % 10000;
    boolean flag = true;
    string label = "rec1676";
};
type R1677 record {
    int id;
    int a = C13 + 1677;
    int b = C40 * 5;
    int c = (1677 + C13) % 10000;
    boolean flag = false;
    string label = "rec1677";
};
type R1678 record {
    int id;
    int a = C14 + 1678;
    int b = C43 * 6;
    int c = (1678 + C14) % 10000;
    boolean flag = true;
    string label = "rec1678";
};
type R1679 record {
    int id;
    int a = C15 + 1679;
    int b = C46 * 7;
    int c = (1679 + C15) % 10000;
    boolean flag = false;
    string label = "rec1679";
};
type R1680 record {
    int id;
    int a = C16 + 1680;
    int b = C49 * 1;
    int c = (1680 + C16) % 10000;
    boolean flag = true;
    string label = "rec1680";
};
type R1681 record {
    int id;
    int a = C17 + 1681;
    int b = C52 * 2;
    int c = (1681 + C17) % 10000;
    boolean flag = false;
    string label = "rec1681";
};
type R1682 record {
    int id;
    int a = C18 + 1682;
    int b = C55 * 3;
    int c = (1682 + C18) % 10000;
    boolean flag = true;
    string label = "rec1682";
};
type R1683 record {
    int id;
    int a = C19 + 1683;
    int b = C58 * 4;
    int c = (1683 + C19) % 10000;
    boolean flag = false;
    string label = "rec1683";
};
type R1684 record {
    int id;
    int a = C20 + 1684;
    int b = C61 * 5;
    int c = (1684 + C20) % 10000;
    boolean flag = true;
    string label = "rec1684";
};
type R1685 record {
    int id;
    int a = C21 + 1685;
    int b = C0 * 6;
    int c = (1685 + C21) % 10000;
    boolean flag = false;
    string label = "rec1685";
};
type R1686 record {
    int id;
    int a = C22 + 1686;
    int b = C3 * 7;
    int c = (1686 + C22) % 10000;
    boolean flag = true;
    string label = "rec1686";
};
type R1687 record {
    int id;
    int a = C23 + 1687;
    int b = C6 * 1;
    int c = (1687 + C23) % 10000;
    boolean flag = false;
    string label = "rec1687";
};
type R1688 record {
    int id;
    int a = C24 + 1688;
    int b = C9 * 2;
    int c = (1688 + C24) % 10000;
    boolean flag = true;
    string label = "rec1688";
};
type R1689 record {
    int id;
    int a = C25 + 1689;
    int b = C12 * 3;
    int c = (1689 + C25) % 10000;
    boolean flag = false;
    string label = "rec1689";
};
type R1690 record {
    int id;
    int a = C26 + 1690;
    int b = C15 * 4;
    int c = (1690 + C26) % 10000;
    boolean flag = true;
    string label = "rec1690";
};
type R1691 record {
    int id;
    int a = C27 + 1691;
    int b = C18 * 5;
    int c = (1691 + C27) % 10000;
    boolean flag = false;
    string label = "rec1691";
};
type R1692 record {
    int id;
    int a = C28 + 1692;
    int b = C21 * 6;
    int c = (1692 + C28) % 10000;
    boolean flag = true;
    string label = "rec1692";
};
type R1693 record {
    int id;
    int a = C29 + 1693;
    int b = C24 * 7;
    int c = (1693 + C29) % 10000;
    boolean flag = false;
    string label = "rec1693";
};
type R1694 record {
    int id;
    int a = C30 + 1694;
    int b = C27 * 1;
    int c = (1694 + C30) % 10000;
    boolean flag = true;
    string label = "rec1694";
};
type R1695 record {
    int id;
    int a = C31 + 1695;
    int b = C30 * 2;
    int c = (1695 + C31) % 10000;
    boolean flag = false;
    string label = "rec1695";
};
type R1696 record {
    int id;
    int a = C32 + 1696;
    int b = C33 * 3;
    int c = (1696 + C32) % 10000;
    boolean flag = true;
    string label = "rec1696";
};
type R1697 record {
    int id;
    int a = C33 + 1697;
    int b = C36 * 4;
    int c = (1697 + C33) % 10000;
    boolean flag = false;
    string label = "rec1697";
};
type R1698 record {
    int id;
    int a = C34 + 1698;
    int b = C39 * 5;
    int c = (1698 + C34) % 10000;
    boolean flag = true;
    string label = "rec1698";
};
type R1699 record {
    int id;
    int a = C35 + 1699;
    int b = C42 * 6;
    int c = (1699 + C35) % 10000;
    boolean flag = false;
    string label = "rec1699";
};
type R1700 record {
    int id;
    int a = C36 + 1700;
    int b = C45 * 7;
    int c = (1700 + C36) % 10000;
    boolean flag = true;
    string label = "rec1700";
};
type R1701 record {
    int id;
    int a = C37 + 1701;
    int b = C48 * 1;
    int c = (1701 + C37) % 10000;
    boolean flag = false;
    string label = "rec1701";
};
type R1702 record {
    int id;
    int a = C38 + 1702;
    int b = C51 * 2;
    int c = (1702 + C38) % 10000;
    boolean flag = true;
    string label = "rec1702";
};
type R1703 record {
    int id;
    int a = C39 + 1703;
    int b = C54 * 3;
    int c = (1703 + C39) % 10000;
    boolean flag = false;
    string label = "rec1703";
};
type R1704 record {
    int id;
    int a = C40 + 1704;
    int b = C57 * 4;
    int c = (1704 + C40) % 10000;
    boolean flag = true;
    string label = "rec1704";
};
type R1705 record {
    int id;
    int a = C41 + 1705;
    int b = C60 * 5;
    int c = (1705 + C41) % 10000;
    boolean flag = false;
    string label = "rec1705";
};
type R1706 record {
    int id;
    int a = C42 + 1706;
    int b = C63 * 6;
    int c = (1706 + C42) % 10000;
    boolean flag = true;
    string label = "rec1706";
};
type R1707 record {
    int id;
    int a = C43 + 1707;
    int b = C2 * 7;
    int c = (1707 + C43) % 10000;
    boolean flag = false;
    string label = "rec1707";
};
type R1708 record {
    int id;
    int a = C44 + 1708;
    int b = C5 * 1;
    int c = (1708 + C44) % 10000;
    boolean flag = true;
    string label = "rec1708";
};
type R1709 record {
    int id;
    int a = C45 + 1709;
    int b = C8 * 2;
    int c = (1709 + C45) % 10000;
    boolean flag = false;
    string label = "rec1709";
};
type R1710 record {
    int id;
    int a = C46 + 1710;
    int b = C11 * 3;
    int c = (1710 + C46) % 10000;
    boolean flag = true;
    string label = "rec1710";
};
type R1711 record {
    int id;
    int a = C47 + 1711;
    int b = C14 * 4;
    int c = (1711 + C47) % 10000;
    boolean flag = false;
    string label = "rec1711";
};
type R1712 record {
    int id;
    int a = C48 + 1712;
    int b = C17 * 5;
    int c = (1712 + C48) % 10000;
    boolean flag = true;
    string label = "rec1712";
};
type R1713 record {
    int id;
    int a = C49 + 1713;
    int b = C20 * 6;
    int c = (1713 + C49) % 10000;
    boolean flag = false;
    string label = "rec1713";
};
type R1714 record {
    int id;
    int a = C50 + 1714;
    int b = C23 * 7;
    int c = (1714 + C50) % 10000;
    boolean flag = true;
    string label = "rec1714";
};
type R1715 record {
    int id;
    int a = C51 + 1715;
    int b = C26 * 1;
    int c = (1715 + C51) % 10000;
    boolean flag = false;
    string label = "rec1715";
};
type R1716 record {
    int id;
    int a = C52 + 1716;
    int b = C29 * 2;
    int c = (1716 + C52) % 10000;
    boolean flag = true;
    string label = "rec1716";
};
type R1717 record {
    int id;
    int a = C53 + 1717;
    int b = C32 * 3;
    int c = (1717 + C53) % 10000;
    boolean flag = false;
    string label = "rec1717";
};
type R1718 record {
    int id;
    int a = C54 + 1718;
    int b = C35 * 4;
    int c = (1718 + C54) % 10000;
    boolean flag = true;
    string label = "rec1718";
};
type R1719 record {
    int id;
    int a = C55 + 1719;
    int b = C38 * 5;
    int c = (1719 + C55) % 10000;
    boolean flag = false;
    string label = "rec1719";
};
type R1720 record {
    int id;
    int a = C56 + 1720;
    int b = C41 * 6;
    int c = (1720 + C56) % 10000;
    boolean flag = true;
    string label = "rec1720";
};
type R1721 record {
    int id;
    int a = C57 + 1721;
    int b = C44 * 7;
    int c = (1721 + C57) % 10000;
    boolean flag = false;
    string label = "rec1721";
};
type R1722 record {
    int id;
    int a = C58 + 1722;
    int b = C47 * 1;
    int c = (1722 + C58) % 10000;
    boolean flag = true;
    string label = "rec1722";
};
type R1723 record {
    int id;
    int a = C59 + 1723;
    int b = C50 * 2;
    int c = (1723 + C59) % 10000;
    boolean flag = false;
    string label = "rec1723";
};
type R1724 record {
    int id;
    int a = C60 + 1724;
    int b = C53 * 3;
    int c = (1724 + C60) % 10000;
    boolean flag = true;
    string label = "rec1724";
};
type R1725 record {
    int id;
    int a = C61 + 1725;
    int b = C56 * 4;
    int c = (1725 + C61) % 10000;
    boolean flag = false;
    string label = "rec1725";
};
type R1726 record {
    int id;
    int a = C62 + 1726;
    int b = C59 * 5;
    int c = (1726 + C62) % 10000;
    boolean flag = true;
    string label = "rec1726";
};
type R1727 record {
    int id;
    int a = C63 + 1727;
    int b = C62 * 6;
    int c = (1727 + C63) % 10000;
    boolean flag = false;
    string label = "rec1727";
};
type R1728 record {
    int id;
    int a = C0 + 1728;
    int b = C1 * 7;
    int c = (1728 + C0) % 10000;
    boolean flag = true;
    string label = "rec1728";
};
type R1729 record {
    int id;
    int a = C1 + 1729;
    int b = C4 * 1;
    int c = (1729 + C1) % 10000;
    boolean flag = false;
    string label = "rec1729";
};
type R1730 record {
    int id;
    int a = C2 + 1730;
    int b = C7 * 2;
    int c = (1730 + C2) % 10000;
    boolean flag = true;
    string label = "rec1730";
};
type R1731 record {
    int id;
    int a = C3 + 1731;
    int b = C10 * 3;
    int c = (1731 + C3) % 10000;
    boolean flag = false;
    string label = "rec1731";
};
type R1732 record {
    int id;
    int a = C4 + 1732;
    int b = C13 * 4;
    int c = (1732 + C4) % 10000;
    boolean flag = true;
    string label = "rec1732";
};
type R1733 record {
    int id;
    int a = C5 + 1733;
    int b = C16 * 5;
    int c = (1733 + C5) % 10000;
    boolean flag = false;
    string label = "rec1733";
};
type R1734 record {
    int id;
    int a = C6 + 1734;
    int b = C19 * 6;
    int c = (1734 + C6) % 10000;
    boolean flag = true;
    string label = "rec1734";
};
type R1735 record {
    int id;
    int a = C7 + 1735;
    int b = C22 * 7;
    int c = (1735 + C7) % 10000;
    boolean flag = false;
    string label = "rec1735";
};
type R1736 record {
    int id;
    int a = C8 + 1736;
    int b = C25 * 1;
    int c = (1736 + C8) % 10000;
    boolean flag = true;
    string label = "rec1736";
};
type R1737 record {
    int id;
    int a = C9 + 1737;
    int b = C28 * 2;
    int c = (1737 + C9) % 10000;
    boolean flag = false;
    string label = "rec1737";
};
type R1738 record {
    int id;
    int a = C10 + 1738;
    int b = C31 * 3;
    int c = (1738 + C10) % 10000;
    boolean flag = true;
    string label = "rec1738";
};
type R1739 record {
    int id;
    int a = C11 + 1739;
    int b = C34 * 4;
    int c = (1739 + C11) % 10000;
    boolean flag = false;
    string label = "rec1739";
};
type R1740 record {
    int id;
    int a = C12 + 1740;
    int b = C37 * 5;
    int c = (1740 + C12) % 10000;
    boolean flag = true;
    string label = "rec1740";
};
type R1741 record {
    int id;
    int a = C13 + 1741;
    int b = C40 * 6;
    int c = (1741 + C13) % 10000;
    boolean flag = false;
    string label = "rec1741";
};
type R1742 record {
    int id;
    int a = C14 + 1742;
    int b = C43 * 7;
    int c = (1742 + C14) % 10000;
    boolean flag = true;
    string label = "rec1742";
};
type R1743 record {
    int id;
    int a = C15 + 1743;
    int b = C46 * 1;
    int c = (1743 + C15) % 10000;
    boolean flag = false;
    string label = "rec1743";
};
type R1744 record {
    int id;
    int a = C16 + 1744;
    int b = C49 * 2;
    int c = (1744 + C16) % 10000;
    boolean flag = true;
    string label = "rec1744";
};
type R1745 record {
    int id;
    int a = C17 + 1745;
    int b = C52 * 3;
    int c = (1745 + C17) % 10000;
    boolean flag = false;
    string label = "rec1745";
};
type R1746 record {
    int id;
    int a = C18 + 1746;
    int b = C55 * 4;
    int c = (1746 + C18) % 10000;
    boolean flag = true;
    string label = "rec1746";
};
type R1747 record {
    int id;
    int a = C19 + 1747;
    int b = C58 * 5;
    int c = (1747 + C19) % 10000;
    boolean flag = false;
    string label = "rec1747";
};
type R1748 record {
    int id;
    int a = C20 + 1748;
    int b = C61 * 6;
    int c = (1748 + C20) % 10000;
    boolean flag = true;
    string label = "rec1748";
};
type R1749 record {
    int id;
    int a = C21 + 1749;
    int b = C0 * 7;
    int c = (1749 + C21) % 10000;
    boolean flag = false;
    string label = "rec1749";
};
type R1750 record {
    int id;
    int a = C22 + 1750;
    int b = C3 * 1;
    int c = (1750 + C22) % 10000;
    boolean flag = true;
    string label = "rec1750";
};
type R1751 record {
    int id;
    int a = C23 + 1751;
    int b = C6 * 2;
    int c = (1751 + C23) % 10000;
    boolean flag = false;
    string label = "rec1751";
};
type R1752 record {
    int id;
    int a = C24 + 1752;
    int b = C9 * 3;
    int c = (1752 + C24) % 10000;
    boolean flag = true;
    string label = "rec1752";
};
type R1753 record {
    int id;
    int a = C25 + 1753;
    int b = C12 * 4;
    int c = (1753 + C25) % 10000;
    boolean flag = false;
    string label = "rec1753";
};
type R1754 record {
    int id;
    int a = C26 + 1754;
    int b = C15 * 5;
    int c = (1754 + C26) % 10000;
    boolean flag = true;
    string label = "rec1754";
};
type R1755 record {
    int id;
    int a = C27 + 1755;
    int b = C18 * 6;
    int c = (1755 + C27) % 10000;
    boolean flag = false;
    string label = "rec1755";
};
type R1756 record {
    int id;
    int a = C28 + 1756;
    int b = C21 * 7;
    int c = (1756 + C28) % 10000;
    boolean flag = true;
    string label = "rec1756";
};
type R1757 record {
    int id;
    int a = C29 + 1757;
    int b = C24 * 1;
    int c = (1757 + C29) % 10000;
    boolean flag = false;
    string label = "rec1757";
};
type R1758 record {
    int id;
    int a = C30 + 1758;
    int b = C27 * 2;
    int c = (1758 + C30) % 10000;
    boolean flag = true;
    string label = "rec1758";
};
type R1759 record {
    int id;
    int a = C31 + 1759;
    int b = C30 * 3;
    int c = (1759 + C31) % 10000;
    boolean flag = false;
    string label = "rec1759";
};
type R1760 record {
    int id;
    int a = C32 + 1760;
    int b = C33 * 4;
    int c = (1760 + C32) % 10000;
    boolean flag = true;
    string label = "rec1760";
};
type R1761 record {
    int id;
    int a = C33 + 1761;
    int b = C36 * 5;
    int c = (1761 + C33) % 10000;
    boolean flag = false;
    string label = "rec1761";
};
type R1762 record {
    int id;
    int a = C34 + 1762;
    int b = C39 * 6;
    int c = (1762 + C34) % 10000;
    boolean flag = true;
    string label = "rec1762";
};
type R1763 record {
    int id;
    int a = C35 + 1763;
    int b = C42 * 7;
    int c = (1763 + C35) % 10000;
    boolean flag = false;
    string label = "rec1763";
};
type R1764 record {
    int id;
    int a = C36 + 1764;
    int b = C45 * 1;
    int c = (1764 + C36) % 10000;
    boolean flag = true;
    string label = "rec1764";
};
type R1765 record {
    int id;
    int a = C37 + 1765;
    int b = C48 * 2;
    int c = (1765 + C37) % 10000;
    boolean flag = false;
    string label = "rec1765";
};
type R1766 record {
    int id;
    int a = C38 + 1766;
    int b = C51 * 3;
    int c = (1766 + C38) % 10000;
    boolean flag = true;
    string label = "rec1766";
};
type R1767 record {
    int id;
    int a = C39 + 1767;
    int b = C54 * 4;
    int c = (1767 + C39) % 10000;
    boolean flag = false;
    string label = "rec1767";
};
type R1768 record {
    int id;
    int a = C40 + 1768;
    int b = C57 * 5;
    int c = (1768 + C40) % 10000;
    boolean flag = true;
    string label = "rec1768";
};
type R1769 record {
    int id;
    int a = C41 + 1769;
    int b = C60 * 6;
    int c = (1769 + C41) % 10000;
    boolean flag = false;
    string label = "rec1769";
};
type R1770 record {
    int id;
    int a = C42 + 1770;
    int b = C63 * 7;
    int c = (1770 + C42) % 10000;
    boolean flag = true;
    string label = "rec1770";
};
type R1771 record {
    int id;
    int a = C43 + 1771;
    int b = C2 * 1;
    int c = (1771 + C43) % 10000;
    boolean flag = false;
    string label = "rec1771";
};
type R1772 record {
    int id;
    int a = C44 + 1772;
    int b = C5 * 2;
    int c = (1772 + C44) % 10000;
    boolean flag = true;
    string label = "rec1772";
};
type R1773 record {
    int id;
    int a = C45 + 1773;
    int b = C8 * 3;
    int c = (1773 + C45) % 10000;
    boolean flag = false;
    string label = "rec1773";
};
type R1774 record {
    int id;
    int a = C46 + 1774;
    int b = C11 * 4;
    int c = (1774 + C46) % 10000;
    boolean flag = true;
    string label = "rec1774";
};
type R1775 record {
    int id;
    int a = C47 + 1775;
    int b = C14 * 5;
    int c = (1775 + C47) % 10000;
    boolean flag = false;
    string label = "rec1775";
};
type R1776 record {
    int id;
    int a = C48 + 1776;
    int b = C17 * 6;
    int c = (1776 + C48) % 10000;
    boolean flag = true;
    string label = "rec1776";
};
type R1777 record {
    int id;
    int a = C49 + 1777;
    int b = C20 * 7;
    int c = (1777 + C49) % 10000;
    boolean flag = false;
    string label = "rec1777";
};
type R1778 record {
    int id;
    int a = C50 + 1778;
    int b = C23 * 1;
    int c = (1778 + C50) % 10000;
    boolean flag = true;
    string label = "rec1778";
};
type R1779 record {
    int id;
    int a = C51 + 1779;
    int b = C26 * 2;
    int c = (1779 + C51) % 10000;
    boolean flag = false;
    string label = "rec1779";
};
type R1780 record {
    int id;
    int a = C52 + 1780;
    int b = C29 * 3;
    int c = (1780 + C52) % 10000;
    boolean flag = true;
    string label = "rec1780";
};
type R1781 record {
    int id;
    int a = C53 + 1781;
    int b = C32 * 4;
    int c = (1781 + C53) % 10000;
    boolean flag = false;
    string label = "rec1781";
};
type R1782 record {
    int id;
    int a = C54 + 1782;
    int b = C35 * 5;
    int c = (1782 + C54) % 10000;
    boolean flag = true;
    string label = "rec1782";
};
type R1783 record {
    int id;
    int a = C55 + 1783;
    int b = C38 * 6;
    int c = (1783 + C55) % 10000;
    boolean flag = false;
    string label = "rec1783";
};
type R1784 record {
    int id;
    int a = C56 + 1784;
    int b = C41 * 7;
    int c = (1784 + C56) % 10000;
    boolean flag = true;
    string label = "rec1784";
};
type R1785 record {
    int id;
    int a = C57 + 1785;
    int b = C44 * 1;
    int c = (1785 + C57) % 10000;
    boolean flag = false;
    string label = "rec1785";
};
type R1786 record {
    int id;
    int a = C58 + 1786;
    int b = C47 * 2;
    int c = (1786 + C58) % 10000;
    boolean flag = true;
    string label = "rec1786";
};
type R1787 record {
    int id;
    int a = C59 + 1787;
    int b = C50 * 3;
    int c = (1787 + C59) % 10000;
    boolean flag = false;
    string label = "rec1787";
};
type R1788 record {
    int id;
    int a = C60 + 1788;
    int b = C53 * 4;
    int c = (1788 + C60) % 10000;
    boolean flag = true;
    string label = "rec1788";
};
type R1789 record {
    int id;
    int a = C61 + 1789;
    int b = C56 * 5;
    int c = (1789 + C61) % 10000;
    boolean flag = false;
    string label = "rec1789";
};
type R1790 record {
    int id;
    int a = C62 + 1790;
    int b = C59 * 6;
    int c = (1790 + C62) % 10000;
    boolean flag = true;
    string label = "rec1790";
};
type R1791 record {
    int id;
    int a = C63 + 1791;
    int b = C62 * 7;
    int c = (1791 + C63) % 10000;
    boolean flag = false;
    string label = "rec1791";
};
type R1792 record {
    int id;
    int a = C0 + 1792;
    int b = C1 * 1;
    int c = (1792 + C0) % 10000;
    boolean flag = true;
    string label = "rec1792";
};
type R1793 record {
    int id;
    int a = C1 + 1793;
    int b = C4 * 2;
    int c = (1793 + C1) % 10000;
    boolean flag = false;
    string label = "rec1793";
};
type R1794 record {
    int id;
    int a = C2 + 1794;
    int b = C7 * 3;
    int c = (1794 + C2) % 10000;
    boolean flag = true;
    string label = "rec1794";
};
type R1795 record {
    int id;
    int a = C3 + 1795;
    int b = C10 * 4;
    int c = (1795 + C3) % 10000;
    boolean flag = false;
    string label = "rec1795";
};
type R1796 record {
    int id;
    int a = C4 + 1796;
    int b = C13 * 5;
    int c = (1796 + C4) % 10000;
    boolean flag = true;
    string label = "rec1796";
};
type R1797 record {
    int id;
    int a = C5 + 1797;
    int b = C16 * 6;
    int c = (1797 + C5) % 10000;
    boolean flag = false;
    string label = "rec1797";
};
type R1798 record {
    int id;
    int a = C6 + 1798;
    int b = C19 * 7;
    int c = (1798 + C6) % 10000;
    boolean flag = true;
    string label = "rec1798";
};
type R1799 record {
    int id;
    int a = C7 + 1799;
    int b = C22 * 1;
    int c = (1799 + C7) % 10000;
    boolean flag = false;
    string label = "rec1799";
};
type R1800 record {
    int id;
    int a = C8 + 1800;
    int b = C25 * 2;
    int c = (1800 + C8) % 10000;
    boolean flag = true;
    string label = "rec1800";
};
type R1801 record {
    int id;
    int a = C9 + 1801;
    int b = C28 * 3;
    int c = (1801 + C9) % 10000;
    boolean flag = false;
    string label = "rec1801";
};
type R1802 record {
    int id;
    int a = C10 + 1802;
    int b = C31 * 4;
    int c = (1802 + C10) % 10000;
    boolean flag = true;
    string label = "rec1802";
};
type R1803 record {
    int id;
    int a = C11 + 1803;
    int b = C34 * 5;
    int c = (1803 + C11) % 10000;
    boolean flag = false;
    string label = "rec1803";
};
type R1804 record {
    int id;
    int a = C12 + 1804;
    int b = C37 * 6;
    int c = (1804 + C12) % 10000;
    boolean flag = true;
    string label = "rec1804";
};
type R1805 record {
    int id;
    int a = C13 + 1805;
    int b = C40 * 7;
    int c = (1805 + C13) % 10000;
    boolean flag = false;
    string label = "rec1805";
};
type R1806 record {
    int id;
    int a = C14 + 1806;
    int b = C43 * 1;
    int c = (1806 + C14) % 10000;
    boolean flag = true;
    string label = "rec1806";
};
type R1807 record {
    int id;
    int a = C15 + 1807;
    int b = C46 * 2;
    int c = (1807 + C15) % 10000;
    boolean flag = false;
    string label = "rec1807";
};
type R1808 record {
    int id;
    int a = C16 + 1808;
    int b = C49 * 3;
    int c = (1808 + C16) % 10000;
    boolean flag = true;
    string label = "rec1808";
};
type R1809 record {
    int id;
    int a = C17 + 1809;
    int b = C52 * 4;
    int c = (1809 + C17) % 10000;
    boolean flag = false;
    string label = "rec1809";
};
type R1810 record {
    int id;
    int a = C18 + 1810;
    int b = C55 * 5;
    int c = (1810 + C18) % 10000;
    boolean flag = true;
    string label = "rec1810";
};
type R1811 record {
    int id;
    int a = C19 + 1811;
    int b = C58 * 6;
    int c = (1811 + C19) % 10000;
    boolean flag = false;
    string label = "rec1811";
};
type R1812 record {
    int id;
    int a = C20 + 1812;
    int b = C61 * 7;
    int c = (1812 + C20) % 10000;
    boolean flag = true;
    string label = "rec1812";
};
type R1813 record {
    int id;
    int a = C21 + 1813;
    int b = C0 * 1;
    int c = (1813 + C21) % 10000;
    boolean flag = false;
    string label = "rec1813";
};
type R1814 record {
    int id;
    int a = C22 + 1814;
    int b = C3 * 2;
    int c = (1814 + C22) % 10000;
    boolean flag = true;
    string label = "rec1814";
};
type R1815 record {
    int id;
    int a = C23 + 1815;
    int b = C6 * 3;
    int c = (1815 + C23) % 10000;
    boolean flag = false;
    string label = "rec1815";
};
type R1816 record {
    int id;
    int a = C24 + 1816;
    int b = C9 * 4;
    int c = (1816 + C24) % 10000;
    boolean flag = true;
    string label = "rec1816";
};
type R1817 record {
    int id;
    int a = C25 + 1817;
    int b = C12 * 5;
    int c = (1817 + C25) % 10000;
    boolean flag = false;
    string label = "rec1817";
};
type R1818 record {
    int id;
    int a = C26 + 1818;
    int b = C15 * 6;
    int c = (1818 + C26) % 10000;
    boolean flag = true;
    string label = "rec1818";
};
type R1819 record {
    int id;
    int a = C27 + 1819;
    int b = C18 * 7;
    int c = (1819 + C27) % 10000;
    boolean flag = false;
    string label = "rec1819";
};
type R1820 record {
    int id;
    int a = C28 + 1820;
    int b = C21 * 1;
    int c = (1820 + C28) % 10000;
    boolean flag = true;
    string label = "rec1820";
};
type R1821 record {
    int id;
    int a = C29 + 1821;
    int b = C24 * 2;
    int c = (1821 + C29) % 10000;
    boolean flag = false;
    string label = "rec1821";
};
type R1822 record {
    int id;
    int a = C30 + 1822;
    int b = C27 * 3;
    int c = (1822 + C30) % 10000;
    boolean flag = true;
    string label = "rec1822";
};
type R1823 record {
    int id;
    int a = C31 + 1823;
    int b = C30 * 4;
    int c = (1823 + C31) % 10000;
    boolean flag = false;
    string label = "rec1823";
};
type R1824 record {
    int id;
    int a = C32 + 1824;
    int b = C33 * 5;
    int c = (1824 + C32) % 10000;
    boolean flag = true;
    string label = "rec1824";
};
type R1825 record {
    int id;
    int a = C33 + 1825;
    int b = C36 * 6;
    int c = (1825 + C33) % 10000;
    boolean flag = false;
    string label = "rec1825";
};
type R1826 record {
    int id;
    int a = C34 + 1826;
    int b = C39 * 7;
    int c = (1826 + C34) % 10000;
    boolean flag = true;
    string label = "rec1826";
};
type R1827 record {
    int id;
    int a = C35 + 1827;
    int b = C42 * 1;
    int c = (1827 + C35) % 10000;
    boolean flag = false;
    string label = "rec1827";
};
type R1828 record {
    int id;
    int a = C36 + 1828;
    int b = C45 * 2;
    int c = (1828 + C36) % 10000;
    boolean flag = true;
    string label = "rec1828";
};
type R1829 record {
    int id;
    int a = C37 + 1829;
    int b = C48 * 3;
    int c = (1829 + C37) % 10000;
    boolean flag = false;
    string label = "rec1829";
};
type R1830 record {
    int id;
    int a = C38 + 1830;
    int b = C51 * 4;
    int c = (1830 + C38) % 10000;
    boolean flag = true;
    string label = "rec1830";
};
type R1831 record {
    int id;
    int a = C39 + 1831;
    int b = C54 * 5;
    int c = (1831 + C39) % 10000;
    boolean flag = false;
    string label = "rec1831";
};
type R1832 record {
    int id;
    int a = C40 + 1832;
    int b = C57 * 6;
    int c = (1832 + C40) % 10000;
    boolean flag = true;
    string label = "rec1832";
};
type R1833 record {
    int id;
    int a = C41 + 1833;
    int b = C60 * 7;
    int c = (1833 + C41) % 10000;
    boolean flag = false;
    string label = "rec1833";
};
type R1834 record {
    int id;
    int a = C42 + 1834;
    int b = C63 * 1;
    int c = (1834 + C42) % 10000;
    boolean flag = true;
    string label = "rec1834";
};
type R1835 record {
    int id;
    int a = C43 + 1835;
    int b = C2 * 2;
    int c = (1835 + C43) % 10000;
    boolean flag = false;
    string label = "rec1835";
};
type R1836 record {
    int id;
    int a = C44 + 1836;
    int b = C5 * 3;
    int c = (1836 + C44) % 10000;
    boolean flag = true;
    string label = "rec1836";
};
type R1837 record {
    int id;
    int a = C45 + 1837;
    int b = C8 * 4;
    int c = (1837 + C45) % 10000;
    boolean flag = false;
    string label = "rec1837";
};
type R1838 record {
    int id;
    int a = C46 + 1838;
    int b = C11 * 5;
    int c = (1838 + C46) % 10000;
    boolean flag = true;
    string label = "rec1838";
};
type R1839 record {
    int id;
    int a = C47 + 1839;
    int b = C14 * 6;
    int c = (1839 + C47) % 10000;
    boolean flag = false;
    string label = "rec1839";
};
type R1840 record {
    int id;
    int a = C48 + 1840;
    int b = C17 * 7;
    int c = (1840 + C48) % 10000;
    boolean flag = true;
    string label = "rec1840";
};
type R1841 record {
    int id;
    int a = C49 + 1841;
    int b = C20 * 1;
    int c = (1841 + C49) % 10000;
    boolean flag = false;
    string label = "rec1841";
};
type R1842 record {
    int id;
    int a = C50 + 1842;
    int b = C23 * 2;
    int c = (1842 + C50) % 10000;
    boolean flag = true;
    string label = "rec1842";
};
type R1843 record {
    int id;
    int a = C51 + 1843;
    int b = C26 * 3;
    int c = (1843 + C51) % 10000;
    boolean flag = false;
    string label = "rec1843";
};
type R1844 record {
    int id;
    int a = C52 + 1844;
    int b = C29 * 4;
    int c = (1844 + C52) % 10000;
    boolean flag = true;
    string label = "rec1844";
};
type R1845 record {
    int id;
    int a = C53 + 1845;
    int b = C32 * 5;
    int c = (1845 + C53) % 10000;
    boolean flag = false;
    string label = "rec1845";
};
type R1846 record {
    int id;
    int a = C54 + 1846;
    int b = C35 * 6;
    int c = (1846 + C54) % 10000;
    boolean flag = true;
    string label = "rec1846";
};
type R1847 record {
    int id;
    int a = C55 + 1847;
    int b = C38 * 7;
    int c = (1847 + C55) % 10000;
    boolean flag = false;
    string label = "rec1847";
};
type R1848 record {
    int id;
    int a = C56 + 1848;
    int b = C41 * 1;
    int c = (1848 + C56) % 10000;
    boolean flag = true;
    string label = "rec1848";
};
type R1849 record {
    int id;
    int a = C57 + 1849;
    int b = C44 * 2;
    int c = (1849 + C57) % 10000;
    boolean flag = false;
    string label = "rec1849";
};
type R1850 record {
    int id;
    int a = C58 + 1850;
    int b = C47 * 3;
    int c = (1850 + C58) % 10000;
    boolean flag = true;
    string label = "rec1850";
};
type R1851 record {
    int id;
    int a = C59 + 1851;
    int b = C50 * 4;
    int c = (1851 + C59) % 10000;
    boolean flag = false;
    string label = "rec1851";
};
type R1852 record {
    int id;
    int a = C60 + 1852;
    int b = C53 * 5;
    int c = (1852 + C60) % 10000;
    boolean flag = true;
    string label = "rec1852";
};
type R1853 record {
    int id;
    int a = C61 + 1853;
    int b = C56 * 6;
    int c = (1853 + C61) % 10000;
    boolean flag = false;
    string label = "rec1853";
};
type R1854 record {
    int id;
    int a = C62 + 1854;
    int b = C59 * 7;
    int c = (1854 + C62) % 10000;
    boolean flag = true;
    string label = "rec1854";
};
type R1855 record {
    int id;
    int a = C63 + 1855;
    int b = C62 * 1;
    int c = (1855 + C63) % 10000;
    boolean flag = false;
    string label = "rec1855";
};
type R1856 record {
    int id;
    int a = C0 + 1856;
    int b = C1 * 2;
    int c = (1856 + C0) % 10000;
    boolean flag = true;
    string label = "rec1856";
};
type R1857 record {
    int id;
    int a = C1 + 1857;
    int b = C4 * 3;
    int c = (1857 + C1) % 10000;
    boolean flag = false;
    string label = "rec1857";
};
type R1858 record {
    int id;
    int a = C2 + 1858;
    int b = C7 * 4;
    int c = (1858 + C2) % 10000;
    boolean flag = true;
    string label = "rec1858";
};
type R1859 record {
    int id;
    int a = C3 + 1859;
    int b = C10 * 5;
    int c = (1859 + C3) % 10000;
    boolean flag = false;
    string label = "rec1859";
};
type R1860 record {
    int id;
    int a = C4 + 1860;
    int b = C13 * 6;
    int c = (1860 + C4) % 10000;
    boolean flag = true;
    string label = "rec1860";
};
type R1861 record {
    int id;
    int a = C5 + 1861;
    int b = C16 * 7;
    int c = (1861 + C5) % 10000;
    boolean flag = false;
    string label = "rec1861";
};
type R1862 record {
    int id;
    int a = C6 + 1862;
    int b = C19 * 1;
    int c = (1862 + C6) % 10000;
    boolean flag = true;
    string label = "rec1862";
};
type R1863 record {
    int id;
    int a = C7 + 1863;
    int b = C22 * 2;
    int c = (1863 + C7) % 10000;
    boolean flag = false;
    string label = "rec1863";
};
type R1864 record {
    int id;
    int a = C8 + 1864;
    int b = C25 * 3;
    int c = (1864 + C8) % 10000;
    boolean flag = true;
    string label = "rec1864";
};
type R1865 record {
    int id;
    int a = C9 + 1865;
    int b = C28 * 4;
    int c = (1865 + C9) % 10000;
    boolean flag = false;
    string label = "rec1865";
};
type R1866 record {
    int id;
    int a = C10 + 1866;
    int b = C31 * 5;
    int c = (1866 + C10) % 10000;
    boolean flag = true;
    string label = "rec1866";
};
type R1867 record {
    int id;
    int a = C11 + 1867;
    int b = C34 * 6;
    int c = (1867 + C11) % 10000;
    boolean flag = false;
    string label = "rec1867";
};
type R1868 record {
    int id;
    int a = C12 + 1868;
    int b = C37 * 7;
    int c = (1868 + C12) % 10000;
    boolean flag = true;
    string label = "rec1868";
};
type R1869 record {
    int id;
    int a = C13 + 1869;
    int b = C40 * 1;
    int c = (1869 + C13) % 10000;
    boolean flag = false;
    string label = "rec1869";
};
type R1870 record {
    int id;
    int a = C14 + 1870;
    int b = C43 * 2;
    int c = (1870 + C14) % 10000;
    boolean flag = true;
    string label = "rec1870";
};
type R1871 record {
    int id;
    int a = C15 + 1871;
    int b = C46 * 3;
    int c = (1871 + C15) % 10000;
    boolean flag = false;
    string label = "rec1871";
};
type R1872 record {
    int id;
    int a = C16 + 1872;
    int b = C49 * 4;
    int c = (1872 + C16) % 10000;
    boolean flag = true;
    string label = "rec1872";
};
type R1873 record {
    int id;
    int a = C17 + 1873;
    int b = C52 * 5;
    int c = (1873 + C17) % 10000;
    boolean flag = false;
    string label = "rec1873";
};
type R1874 record {
    int id;
    int a = C18 + 1874;
    int b = C55 * 6;
    int c = (1874 + C18) % 10000;
    boolean flag = true;
    string label = "rec1874";
};
type R1875 record {
    int id;
    int a = C19 + 1875;
    int b = C58 * 7;
    int c = (1875 + C19) % 10000;
    boolean flag = false;
    string label = "rec1875";
};
type R1876 record {
    int id;
    int a = C20 + 1876;
    int b = C61 * 1;
    int c = (1876 + C20) % 10000;
    boolean flag = true;
    string label = "rec1876";
};
type R1877 record {
    int id;
    int a = C21 + 1877;
    int b = C0 * 2;
    int c = (1877 + C21) % 10000;
    boolean flag = false;
    string label = "rec1877";
};
type R1878 record {
    int id;
    int a = C22 + 1878;
    int b = C3 * 3;
    int c = (1878 + C22) % 10000;
    boolean flag = true;
    string label = "rec1878";
};
type R1879 record {
    int id;
    int a = C23 + 1879;
    int b = C6 * 4;
    int c = (1879 + C23) % 10000;
    boolean flag = false;
    string label = "rec1879";
};
type R1880 record {
    int id;
    int a = C24 + 1880;
    int b = C9 * 5;
    int c = (1880 + C24) % 10000;
    boolean flag = true;
    string label = "rec1880";
};
type R1881 record {
    int id;
    int a = C25 + 1881;
    int b = C12 * 6;
    int c = (1881 + C25) % 10000;
    boolean flag = false;
    string label = "rec1881";
};
type R1882 record {
    int id;
    int a = C26 + 1882;
    int b = C15 * 7;
    int c = (1882 + C26) % 10000;
    boolean flag = true;
    string label = "rec1882";
};
type R1883 record {
    int id;
    int a = C27 + 1883;
    int b = C18 * 1;
    int c = (1883 + C27) % 10000;
    boolean flag = false;
    string label = "rec1883";
};
type R1884 record {
    int id;
    int a = C28 + 1884;
    int b = C21 * 2;
    int c = (1884 + C28) % 10000;
    boolean flag = true;
    string label = "rec1884";
};
type R1885 record {
    int id;
    int a = C29 + 1885;
    int b = C24 * 3;
    int c = (1885 + C29) % 10000;
    boolean flag = false;
    string label = "rec1885";
};
type R1886 record {
    int id;
    int a = C30 + 1886;
    int b = C27 * 4;
    int c = (1886 + C30) % 10000;
    boolean flag = true;
    string label = "rec1886";
};
type R1887 record {
    int id;
    int a = C31 + 1887;
    int b = C30 * 5;
    int c = (1887 + C31) % 10000;
    boolean flag = false;
    string label = "rec1887";
};
type R1888 record {
    int id;
    int a = C32 + 1888;
    int b = C33 * 6;
    int c = (1888 + C32) % 10000;
    boolean flag = true;
    string label = "rec1888";
};
type R1889 record {
    int id;
    int a = C33 + 1889;
    int b = C36 * 7;
    int c = (1889 + C33) % 10000;
    boolean flag = false;
    string label = "rec1889";
};
type R1890 record {
    int id;
    int a = C34 + 1890;
    int b = C39 * 1;
    int c = (1890 + C34) % 10000;
    boolean flag = true;
    string label = "rec1890";
};
type R1891 record {
    int id;
    int a = C35 + 1891;
    int b = C42 * 2;
    int c = (1891 + C35) % 10000;
    boolean flag = false;
    string label = "rec1891";
};
type R1892 record {
    int id;
    int a = C36 + 1892;
    int b = C45 * 3;
    int c = (1892 + C36) % 10000;
    boolean flag = true;
    string label = "rec1892";
};
type R1893 record {
    int id;
    int a = C37 + 1893;
    int b = C48 * 4;
    int c = (1893 + C37) % 10000;
    boolean flag = false;
    string label = "rec1893";
};
type R1894 record {
    int id;
    int a = C38 + 1894;
    int b = C51 * 5;
    int c = (1894 + C38) % 10000;
    boolean flag = true;
    string label = "rec1894";
};
type R1895 record {
    int id;
    int a = C39 + 1895;
    int b = C54 * 6;
    int c = (1895 + C39) % 10000;
    boolean flag = false;
    string label = "rec1895";
};
type R1896 record {
    int id;
    int a = C40 + 1896;
    int b = C57 * 7;
    int c = (1896 + C40) % 10000;
    boolean flag = true;
    string label = "rec1896";
};
type R1897 record {
    int id;
    int a = C41 + 1897;
    int b = C60 * 1;
    int c = (1897 + C41) % 10000;
    boolean flag = false;
    string label = "rec1897";
};
type R1898 record {
    int id;
    int a = C42 + 1898;
    int b = C63 * 2;
    int c = (1898 + C42) % 10000;
    boolean flag = true;
    string label = "rec1898";
};
type R1899 record {
    int id;
    int a = C43 + 1899;
    int b = C2 * 3;
    int c = (1899 + C43) % 10000;
    boolean flag = false;
    string label = "rec1899";
};
type R1900 record {
    int id;
    int a = C44 + 1900;
    int b = C5 * 4;
    int c = (1900 + C44) % 10000;
    boolean flag = true;
    string label = "rec1900";
};
type R1901 record {
    int id;
    int a = C45 + 1901;
    int b = C8 * 5;
    int c = (1901 + C45) % 10000;
    boolean flag = false;
    string label = "rec1901";
};
type R1902 record {
    int id;
    int a = C46 + 1902;
    int b = C11 * 6;
    int c = (1902 + C46) % 10000;
    boolean flag = true;
    string label = "rec1902";
};
type R1903 record {
    int id;
    int a = C47 + 1903;
    int b = C14 * 7;
    int c = (1903 + C47) % 10000;
    boolean flag = false;
    string label = "rec1903";
};
type R1904 record {
    int id;
    int a = C48 + 1904;
    int b = C17 * 1;
    int c = (1904 + C48) % 10000;
    boolean flag = true;
    string label = "rec1904";
};
type R1905 record {
    int id;
    int a = C49 + 1905;
    int b = C20 * 2;
    int c = (1905 + C49) % 10000;
    boolean flag = false;
    string label = "rec1905";
};
type R1906 record {
    int id;
    int a = C50 + 1906;
    int b = C23 * 3;
    int c = (1906 + C50) % 10000;
    boolean flag = true;
    string label = "rec1906";
};
type R1907 record {
    int id;
    int a = C51 + 1907;
    int b = C26 * 4;
    int c = (1907 + C51) % 10000;
    boolean flag = false;
    string label = "rec1907";
};
type R1908 record {
    int id;
    int a = C52 + 1908;
    int b = C29 * 5;
    int c = (1908 + C52) % 10000;
    boolean flag = true;
    string label = "rec1908";
};
type R1909 record {
    int id;
    int a = C53 + 1909;
    int b = C32 * 6;
    int c = (1909 + C53) % 10000;
    boolean flag = false;
    string label = "rec1909";
};
type R1910 record {
    int id;
    int a = C54 + 1910;
    int b = C35 * 7;
    int c = (1910 + C54) % 10000;
    boolean flag = true;
    string label = "rec1910";
};
type R1911 record {
    int id;
    int a = C55 + 1911;
    int b = C38 * 1;
    int c = (1911 + C55) % 10000;
    boolean flag = false;
    string label = "rec1911";
};
type R1912 record {
    int id;
    int a = C56 + 1912;
    int b = C41 * 2;
    int c = (1912 + C56) % 10000;
    boolean flag = true;
    string label = "rec1912";
};
type R1913 record {
    int id;
    int a = C57 + 1913;
    int b = C44 * 3;
    int c = (1913 + C57) % 10000;
    boolean flag = false;
    string label = "rec1913";
};
type R1914 record {
    int id;
    int a = C58 + 1914;
    int b = C47 * 4;
    int c = (1914 + C58) % 10000;
    boolean flag = true;
    string label = "rec1914";
};
type R1915 record {
    int id;
    int a = C59 + 1915;
    int b = C50 * 5;
    int c = (1915 + C59) % 10000;
    boolean flag = false;
    string label = "rec1915";
};
type R1916 record {
    int id;
    int a = C60 + 1916;
    int b = C53 * 6;
    int c = (1916 + C60) % 10000;
    boolean flag = true;
    string label = "rec1916";
};
type R1917 record {
    int id;
    int a = C61 + 1917;
    int b = C56 * 7;
    int c = (1917 + C61) % 10000;
    boolean flag = false;
    string label = "rec1917";
};
type R1918 record {
    int id;
    int a = C62 + 1918;
    int b = C59 * 1;
    int c = (1918 + C62) % 10000;
    boolean flag = true;
    string label = "rec1918";
};
type R1919 record {
    int id;
    int a = C63 + 1919;
    int b = C62 * 2;
    int c = (1919 + C63) % 10000;
    boolean flag = false;
    string label = "rec1919";
};
type R1920 record {
    int id;
    int a = C0 + 1920;
    int b = C1 * 3;
    int c = (1920 + C0) % 10000;
    boolean flag = true;
    string label = "rec1920";
};
type R1921 record {
    int id;
    int a = C1 + 1921;
    int b = C4 * 4;
    int c = (1921 + C1) % 10000;
    boolean flag = false;
    string label = "rec1921";
};
type R1922 record {
    int id;
    int a = C2 + 1922;
    int b = C7 * 5;
    int c = (1922 + C2) % 10000;
    boolean flag = true;
    string label = "rec1922";
};
type R1923 record {
    int id;
    int a = C3 + 1923;
    int b = C10 * 6;
    int c = (1923 + C3) % 10000;
    boolean flag = false;
    string label = "rec1923";
};
type R1924 record {
    int id;
    int a = C4 + 1924;
    int b = C13 * 7;
    int c = (1924 + C4) % 10000;
    boolean flag = true;
    string label = "rec1924";
};
type R1925 record {
    int id;
    int a = C5 + 1925;
    int b = C16 * 1;
    int c = (1925 + C5) % 10000;
    boolean flag = false;
    string label = "rec1925";
};
type R1926 record {
    int id;
    int a = C6 + 1926;
    int b = C19 * 2;
    int c = (1926 + C6) % 10000;
    boolean flag = true;
    string label = "rec1926";
};
type R1927 record {
    int id;
    int a = C7 + 1927;
    int b = C22 * 3;
    int c = (1927 + C7) % 10000;
    boolean flag = false;
    string label = "rec1927";
};
type R1928 record {
    int id;
    int a = C8 + 1928;
    int b = C25 * 4;
    int c = (1928 + C8) % 10000;
    boolean flag = true;
    string label = "rec1928";
};
type R1929 record {
    int id;
    int a = C9 + 1929;
    int b = C28 * 5;
    int c = (1929 + C9) % 10000;
    boolean flag = false;
    string label = "rec1929";
};
type R1930 record {
    int id;
    int a = C10 + 1930;
    int b = C31 * 6;
    int c = (1930 + C10) % 10000;
    boolean flag = true;
    string label = "rec1930";
};
type R1931 record {
    int id;
    int a = C11 + 1931;
    int b = C34 * 7;
    int c = (1931 + C11) % 10000;
    boolean flag = false;
    string label = "rec1931";
};
type R1932 record {
    int id;
    int a = C12 + 1932;
    int b = C37 * 1;
    int c = (1932 + C12) % 10000;
    boolean flag = true;
    string label = "rec1932";
};
type R1933 record {
    int id;
    int a = C13 + 1933;
    int b = C40 * 2;
    int c = (1933 + C13) % 10000;
    boolean flag = false;
    string label = "rec1933";
};
type R1934 record {
    int id;
    int a = C14 + 1934;
    int b = C43 * 3;
    int c = (1934 + C14) % 10000;
    boolean flag = true;
    string label = "rec1934";
};
type R1935 record {
    int id;
    int a = C15 + 1935;
    int b = C46 * 4;
    int c = (1935 + C15) % 10000;
    boolean flag = false;
    string label = "rec1935";
};
type R1936 record {
    int id;
    int a = C16 + 1936;
    int b = C49 * 5;
    int c = (1936 + C16) % 10000;
    boolean flag = true;
    string label = "rec1936";
};
type R1937 record {
    int id;
    int a = C17 + 1937;
    int b = C52 * 6;
    int c = (1937 + C17) % 10000;
    boolean flag = false;
    string label = "rec1937";
};
type R1938 record {
    int id;
    int a = C18 + 1938;
    int b = C55 * 7;
    int c = (1938 + C18) % 10000;
    boolean flag = true;
    string label = "rec1938";
};
type R1939 record {
    int id;
    int a = C19 + 1939;
    int b = C58 * 1;
    int c = (1939 + C19) % 10000;
    boolean flag = false;
    string label = "rec1939";
};
type R1940 record {
    int id;
    int a = C20 + 1940;
    int b = C61 * 2;
    int c = (1940 + C20) % 10000;
    boolean flag = true;
    string label = "rec1940";
};
type R1941 record {
    int id;
    int a = C21 + 1941;
    int b = C0 * 3;
    int c = (1941 + C21) % 10000;
    boolean flag = false;
    string label = "rec1941";
};
type R1942 record {
    int id;
    int a = C22 + 1942;
    int b = C3 * 4;
    int c = (1942 + C22) % 10000;
    boolean flag = true;
    string label = "rec1942";
};
type R1943 record {
    int id;
    int a = C23 + 1943;
    int b = C6 * 5;
    int c = (1943 + C23) % 10000;
    boolean flag = false;
    string label = "rec1943";
};
type R1944 record {
    int id;
    int a = C24 + 1944;
    int b = C9 * 6;
    int c = (1944 + C24) % 10000;
    boolean flag = true;
    string label = "rec1944";
};
type R1945 record {
    int id;
    int a = C25 + 1945;
    int b = C12 * 7;
    int c = (1945 + C25) % 10000;
    boolean flag = false;
    string label = "rec1945";
};
type R1946 record {
    int id;
    int a = C26 + 1946;
    int b = C15 * 1;
    int c = (1946 + C26) % 10000;
    boolean flag = true;
    string label = "rec1946";
};
type R1947 record {
    int id;
    int a = C27 + 1947;
    int b = C18 * 2;
    int c = (1947 + C27) % 10000;
    boolean flag = false;
    string label = "rec1947";
};
type R1948 record {
    int id;
    int a = C28 + 1948;
    int b = C21 * 3;
    int c = (1948 + C28) % 10000;
    boolean flag = true;
    string label = "rec1948";
};
type R1949 record {
    int id;
    int a = C29 + 1949;
    int b = C24 * 4;
    int c = (1949 + C29) % 10000;
    boolean flag = false;
    string label = "rec1949";
};
type R1950 record {
    int id;
    int a = C30 + 1950;
    int b = C27 * 5;
    int c = (1950 + C30) % 10000;
    boolean flag = true;
    string label = "rec1950";
};
type R1951 record {
    int id;
    int a = C31 + 1951;
    int b = C30 * 6;
    int c = (1951 + C31) % 10000;
    boolean flag = false;
    string label = "rec1951";
};
type R1952 record {
    int id;
    int a = C32 + 1952;
    int b = C33 * 7;
    int c = (1952 + C32) % 10000;
    boolean flag = true;
    string label = "rec1952";
};
type R1953 record {
    int id;
    int a = C33 + 1953;
    int b = C36 * 1;
    int c = (1953 + C33) % 10000;
    boolean flag = false;
    string label = "rec1953";
};
type R1954 record {
    int id;
    int a = C34 + 1954;
    int b = C39 * 2;
    int c = (1954 + C34) % 10000;
    boolean flag = true;
    string label = "rec1954";
};
type R1955 record {
    int id;
    int a = C35 + 1955;
    int b = C42 * 3;
    int c = (1955 + C35) % 10000;
    boolean flag = false;
    string label = "rec1955";
};
type R1956 record {
    int id;
    int a = C36 + 1956;
    int b = C45 * 4;
    int c = (1956 + C36) % 10000;
    boolean flag = true;
    string label = "rec1956";
};
type R1957 record {
    int id;
    int a = C37 + 1957;
    int b = C48 * 5;
    int c = (1957 + C37) % 10000;
    boolean flag = false;
    string label = "rec1957";
};
type R1958 record {
    int id;
    int a = C38 + 1958;
    int b = C51 * 6;
    int c = (1958 + C38) % 10000;
    boolean flag = true;
    string label = "rec1958";
};
type R1959 record {
    int id;
    int a = C39 + 1959;
    int b = C54 * 7;
    int c = (1959 + C39) % 10000;
    boolean flag = false;
    string label = "rec1959";
};
type R1960 record {
    int id;
    int a = C40 + 1960;
    int b = C57 * 1;
    int c = (1960 + C40) % 10000;
    boolean flag = true;
    string label = "rec1960";
};
type R1961 record {
    int id;
    int a = C41 + 1961;
    int b = C60 * 2;
    int c = (1961 + C41) % 10000;
    boolean flag = false;
    string label = "rec1961";
};
type R1962 record {
    int id;
    int a = C42 + 1962;
    int b = C63 * 3;
    int c = (1962 + C42) % 10000;
    boolean flag = true;
    string label = "rec1962";
};
type R1963 record {
    int id;
    int a = C43 + 1963;
    int b = C2 * 4;
    int c = (1963 + C43) % 10000;
    boolean flag = false;
    string label = "rec1963";
};
type R1964 record {
    int id;
    int a = C44 + 1964;
    int b = C5 * 5;
    int c = (1964 + C44) % 10000;
    boolean flag = true;
    string label = "rec1964";
};
type R1965 record {
    int id;
    int a = C45 + 1965;
    int b = C8 * 6;
    int c = (1965 + C45) % 10000;
    boolean flag = false;
    string label = "rec1965";
};
type R1966 record {
    int id;
    int a = C46 + 1966;
    int b = C11 * 7;
    int c = (1966 + C46) % 10000;
    boolean flag = true;
    string label = "rec1966";
};
type R1967 record {
    int id;
    int a = C47 + 1967;
    int b = C14 * 1;
    int c = (1967 + C47) % 10000;
    boolean flag = false;
    string label = "rec1967";
};
type R1968 record {
    int id;
    int a = C48 + 1968;
    int b = C17 * 2;
    int c = (1968 + C48) % 10000;
    boolean flag = true;
    string label = "rec1968";
};
type R1969 record {
    int id;
    int a = C49 + 1969;
    int b = C20 * 3;
    int c = (1969 + C49) % 10000;
    boolean flag = false;
    string label = "rec1969";
};
type R1970 record {
    int id;
    int a = C50 + 1970;
    int b = C23 * 4;
    int c = (1970 + C50) % 10000;
    boolean flag = true;
    string label = "rec1970";
};
type R1971 record {
    int id;
    int a = C51 + 1971;
    int b = C26 * 5;
    int c = (1971 + C51) % 10000;
    boolean flag = false;
    string label = "rec1971";
};
type R1972 record {
    int id;
    int a = C52 + 1972;
    int b = C29 * 6;
    int c = (1972 + C52) % 10000;
    boolean flag = true;
    string label = "rec1972";
};
type R1973 record {
    int id;
    int a = C53 + 1973;
    int b = C32 * 7;
    int c = (1973 + C53) % 10000;
    boolean flag = false;
    string label = "rec1973";
};
type R1974 record {
    int id;
    int a = C54 + 1974;
    int b = C35 * 1;
    int c = (1974 + C54) % 10000;
    boolean flag = true;
    string label = "rec1974";
};
type R1975 record {
    int id;
    int a = C55 + 1975;
    int b = C38 * 2;
    int c = (1975 + C55) % 10000;
    boolean flag = false;
    string label = "rec1975";
};
type R1976 record {
    int id;
    int a = C56 + 1976;
    int b = C41 * 3;
    int c = (1976 + C56) % 10000;
    boolean flag = true;
    string label = "rec1976";
};
type R1977 record {
    int id;
    int a = C57 + 1977;
    int b = C44 * 4;
    int c = (1977 + C57) % 10000;
    boolean flag = false;
    string label = "rec1977";
};
type R1978 record {
    int id;
    int a = C58 + 1978;
    int b = C47 * 5;
    int c = (1978 + C58) % 10000;
    boolean flag = true;
    string label = "rec1978";
};
type R1979 record {
    int id;
    int a = C59 + 1979;
    int b = C50 * 6;
    int c = (1979 + C59) % 10000;
    boolean flag = false;
    string label = "rec1979";
};
type R1980 record {
    int id;
    int a = C60 + 1980;
    int b = C53 * 7;
    int c = (1980 + C60) % 10000;
    boolean flag = true;
    string label = "rec1980";
};
type R1981 record {
    int id;
    int a = C61 + 1981;
    int b = C56 * 1;
    int c = (1981 + C61) % 10000;
    boolean flag = false;
    string label = "rec1981";
};
type R1982 record {
    int id;
    int a = C62 + 1982;
    int b = C59 * 2;
    int c = (1982 + C62) % 10000;
    boolean flag = true;
    string label = "rec1982";
};
type R1983 record {
    int id;
    int a = C63 + 1983;
    int b = C62 * 3;
    int c = (1983 + C63) % 10000;
    boolean flag = false;
    string label = "rec1983";
};
type R1984 record {
    int id;
    int a = C0 + 1984;
    int b = C1 * 4;
    int c = (1984 + C0) % 10000;
    boolean flag = true;
    string label = "rec1984";
};
type R1985 record {
    int id;
    int a = C1 + 1985;
    int b = C4 * 5;
    int c = (1985 + C1) % 10000;
    boolean flag = false;
    string label = "rec1985";
};
type R1986 record {
    int id;
    int a = C2 + 1986;
    int b = C7 * 6;
    int c = (1986 + C2) % 10000;
    boolean flag = true;
    string label = "rec1986";
};
type R1987 record {
    int id;
    int a = C3 + 1987;
    int b = C10 * 7;
    int c = (1987 + C3) % 10000;
    boolean flag = false;
    string label = "rec1987";
};
type R1988 record {
    int id;
    int a = C4 + 1988;
    int b = C13 * 1;
    int c = (1988 + C4) % 10000;
    boolean flag = true;
    string label = "rec1988";
};
type R1989 record {
    int id;
    int a = C5 + 1989;
    int b = C16 * 2;
    int c = (1989 + C5) % 10000;
    boolean flag = false;
    string label = "rec1989";
};
type R1990 record {
    int id;
    int a = C6 + 1990;
    int b = C19 * 3;
    int c = (1990 + C6) % 10000;
    boolean flag = true;
    string label = "rec1990";
};
type R1991 record {
    int id;
    int a = C7 + 1991;
    int b = C22 * 4;
    int c = (1991 + C7) % 10000;
    boolean flag = false;
    string label = "rec1991";
};
type R1992 record {
    int id;
    int a = C8 + 1992;
    int b = C25 * 5;
    int c = (1992 + C8) % 10000;
    boolean flag = true;
    string label = "rec1992";
};
type R1993 record {
    int id;
    int a = C9 + 1993;
    int b = C28 * 6;
    int c = (1993 + C9) % 10000;
    boolean flag = false;
    string label = "rec1993";
};
type R1994 record {
    int id;
    int a = C10 + 1994;
    int b = C31 * 7;
    int c = (1994 + C10) % 10000;
    boolean flag = true;
    string label = "rec1994";
};
type R1995 record {
    int id;
    int a = C11 + 1995;
    int b = C34 * 1;
    int c = (1995 + C11) % 10000;
    boolean flag = false;
    string label = "rec1995";
};
type R1996 record {
    int id;
    int a = C12 + 1996;
    int b = C37 * 2;
    int c = (1996 + C12) % 10000;
    boolean flag = true;
    string label = "rec1996";
};
type R1997 record {
    int id;
    int a = C13 + 1997;
    int b = C40 * 3;
    int c = (1997 + C13) % 10000;
    boolean flag = false;
    string label = "rec1997";
};
type R1998 record {
    int id;
    int a = C14 + 1998;
    int b = C43 * 4;
    int c = (1998 + C14) % 10000;
    boolean flag = true;
    string label = "rec1998";
};
type R1999 record {
    int id;
    int a = C15 + 1999;
    int b = C46 * 5;
    int c = (1999 + C15) % 10000;
    boolean flag = false;
    string label = "rec1999";
};
type R2000 record {
    int id;
    int a = C16 + 2000;
    int b = C49 * 6;
    int c = (2000 + C16) % 10000;
    boolean flag = true;
    string label = "rec2000";
};
type R2001 record {
    int id;
    int a = C17 + 2001;
    int b = C52 * 7;
    int c = (2001 + C17) % 10000;
    boolean flag = false;
    string label = "rec2001";
};
type R2002 record {
    int id;
    int a = C18 + 2002;
    int b = C55 * 1;
    int c = (2002 + C18) % 10000;
    boolean flag = true;
    string label = "rec2002";
};
type R2003 record {
    int id;
    int a = C19 + 2003;
    int b = C58 * 2;
    int c = (2003 + C19) % 10000;
    boolean flag = false;
    string label = "rec2003";
};
type R2004 record {
    int id;
    int a = C20 + 2004;
    int b = C61 * 3;
    int c = (2004 + C20) % 10000;
    boolean flag = true;
    string label = "rec2004";
};
type R2005 record {
    int id;
    int a = C21 + 2005;
    int b = C0 * 4;
    int c = (2005 + C21) % 10000;
    boolean flag = false;
    string label = "rec2005";
};
type R2006 record {
    int id;
    int a = C22 + 2006;
    int b = C3 * 5;
    int c = (2006 + C22) % 10000;
    boolean flag = true;
    string label = "rec2006";
};
type R2007 record {
    int id;
    int a = C23 + 2007;
    int b = C6 * 6;
    int c = (2007 + C23) % 10000;
    boolean flag = false;
    string label = "rec2007";
};
type R2008 record {
    int id;
    int a = C24 + 2008;
    int b = C9 * 7;
    int c = (2008 + C24) % 10000;
    boolean flag = true;
    string label = "rec2008";
};
type R2009 record {
    int id;
    int a = C25 + 2009;
    int b = C12 * 1;
    int c = (2009 + C25) % 10000;
    boolean flag = false;
    string label = "rec2009";
};
type R2010 record {
    int id;
    int a = C26 + 2010;
    int b = C15 * 2;
    int c = (2010 + C26) % 10000;
    boolean flag = true;
    string label = "rec2010";
};
type R2011 record {
    int id;
    int a = C27 + 2011;
    int b = C18 * 3;
    int c = (2011 + C27) % 10000;
    boolean flag = false;
    string label = "rec2011";
};
type R2012 record {
    int id;
    int a = C28 + 2012;
    int b = C21 * 4;
    int c = (2012 + C28) % 10000;
    boolean flag = true;
    string label = "rec2012";
};
type R2013 record {
    int id;
    int a = C29 + 2013;
    int b = C24 * 5;
    int c = (2013 + C29) % 10000;
    boolean flag = false;
    string label = "rec2013";
};
type R2014 record {
    int id;
    int a = C30 + 2014;
    int b = C27 * 6;
    int c = (2014 + C30) % 10000;
    boolean flag = true;
    string label = "rec2014";
};
type R2015 record {
    int id;
    int a = C31 + 2015;
    int b = C30 * 7;
    int c = (2015 + C31) % 10000;
    boolean flag = false;
    string label = "rec2015";
};
type R2016 record {
    int id;
    int a = C32 + 2016;
    int b = C33 * 1;
    int c = (2016 + C32) % 10000;
    boolean flag = true;
    string label = "rec2016";
};
type R2017 record {
    int id;
    int a = C33 + 2017;
    int b = C36 * 2;
    int c = (2017 + C33) % 10000;
    boolean flag = false;
    string label = "rec2017";
};
type R2018 record {
    int id;
    int a = C34 + 2018;
    int b = C39 * 3;
    int c = (2018 + C34) % 10000;
    boolean flag = true;
    string label = "rec2018";
};
type R2019 record {
    int id;
    int a = C35 + 2019;
    int b = C42 * 4;
    int c = (2019 + C35) % 10000;
    boolean flag = false;
    string label = "rec2019";
};
type R2020 record {
    int id;
    int a = C36 + 2020;
    int b = C45 * 5;
    int c = (2020 + C36) % 10000;
    boolean flag = true;
    string label = "rec2020";
};
type R2021 record {
    int id;
    int a = C37 + 2021;
    int b = C48 * 6;
    int c = (2021 + C37) % 10000;
    boolean flag = false;
    string label = "rec2021";
};
type R2022 record {
    int id;
    int a = C38 + 2022;
    int b = C51 * 7;
    int c = (2022 + C38) % 10000;
    boolean flag = true;
    string label = "rec2022";
};
type R2023 record {
    int id;
    int a = C39 + 2023;
    int b = C54 * 1;
    int c = (2023 + C39) % 10000;
    boolean flag = false;
    string label = "rec2023";
};
type R2024 record {
    int id;
    int a = C40 + 2024;
    int b = C57 * 2;
    int c = (2024 + C40) % 10000;
    boolean flag = true;
    string label = "rec2024";
};
type R2025 record {
    int id;
    int a = C41 + 2025;
    int b = C60 * 3;
    int c = (2025 + C41) % 10000;
    boolean flag = false;
    string label = "rec2025";
};
type R2026 record {
    int id;
    int a = C42 + 2026;
    int b = C63 * 4;
    int c = (2026 + C42) % 10000;
    boolean flag = true;
    string label = "rec2026";
};
type R2027 record {
    int id;
    int a = C43 + 2027;
    int b = C2 * 5;
    int c = (2027 + C43) % 10000;
    boolean flag = false;
    string label = "rec2027";
};
type R2028 record {
    int id;
    int a = C44 + 2028;
    int b = C5 * 6;
    int c = (2028 + C44) % 10000;
    boolean flag = true;
    string label = "rec2028";
};
type R2029 record {
    int id;
    int a = C45 + 2029;
    int b = C8 * 7;
    int c = (2029 + C45) % 10000;
    boolean flag = false;
    string label = "rec2029";
};
type R2030 record {
    int id;
    int a = C46 + 2030;
    int b = C11 * 1;
    int c = (2030 + C46) % 10000;
    boolean flag = true;
    string label = "rec2030";
};
type R2031 record {
    int id;
    int a = C47 + 2031;
    int b = C14 * 2;
    int c = (2031 + C47) % 10000;
    boolean flag = false;
    string label = "rec2031";
};
type R2032 record {
    int id;
    int a = C48 + 2032;
    int b = C17 * 3;
    int c = (2032 + C48) % 10000;
    boolean flag = true;
    string label = "rec2032";
};
type R2033 record {
    int id;
    int a = C49 + 2033;
    int b = C20 * 4;
    int c = (2033 + C49) % 10000;
    boolean flag = false;
    string label = "rec2033";
};
type R2034 record {
    int id;
    int a = C50 + 2034;
    int b = C23 * 5;
    int c = (2034 + C50) % 10000;
    boolean flag = true;
    string label = "rec2034";
};
type R2035 record {
    int id;
    int a = C51 + 2035;
    int b = C26 * 6;
    int c = (2035 + C51) % 10000;
    boolean flag = false;
    string label = "rec2035";
};
type R2036 record {
    int id;
    int a = C52 + 2036;
    int b = C29 * 7;
    int c = (2036 + C52) % 10000;
    boolean flag = true;
    string label = "rec2036";
};
type R2037 record {
    int id;
    int a = C53 + 2037;
    int b = C32 * 1;
    int c = (2037 + C53) % 10000;
    boolean flag = false;
    string label = "rec2037";
};
type R2038 record {
    int id;
    int a = C54 + 2038;
    int b = C35 * 2;
    int c = (2038 + C54) % 10000;
    boolean flag = true;
    string label = "rec2038";
};
type R2039 record {
    int id;
    int a = C55 + 2039;
    int b = C38 * 3;
    int c = (2039 + C55) % 10000;
    boolean flag = false;
    string label = "rec2039";
};
type R2040 record {
    int id;
    int a = C56 + 2040;
    int b = C41 * 4;
    int c = (2040 + C56) % 10000;
    boolean flag = true;
    string label = "rec2040";
};
type R2041 record {
    int id;
    int a = C57 + 2041;
    int b = C44 * 5;
    int c = (2041 + C57) % 10000;
    boolean flag = false;
    string label = "rec2041";
};
type R2042 record {
    int id;
    int a = C58 + 2042;
    int b = C47 * 6;
    int c = (2042 + C58) % 10000;
    boolean flag = true;
    string label = "rec2042";
};
type R2043 record {
    int id;
    int a = C59 + 2043;
    int b = C50 * 7;
    int c = (2043 + C59) % 10000;
    boolean flag = false;
    string label = "rec2043";
};
type R2044 record {
    int id;
    int a = C60 + 2044;
    int b = C53 * 1;
    int c = (2044 + C60) % 10000;
    boolean flag = true;
    string label = "rec2044";
};
type R2045 record {
    int id;
    int a = C61 + 2045;
    int b = C56 * 2;
    int c = (2045 + C61) % 10000;
    boolean flag = false;
    string label = "rec2045";
};
type R2046 record {
    int id;
    int a = C62 + 2046;
    int b = C59 * 3;
    int c = (2046 + C62) % 10000;
    boolean flag = true;
    string label = "rec2046";
};
type R2047 record {
    int id;
    int a = C63 + 2047;
    int b = C62 * 4;
    int c = (2047 + C63) % 10000;
    boolean flag = false;
    string label = "rec2047";
};
type R2048 record {
    int id;
    int a = C0 + 2048;
    int b = C1 * 5;
    int c = (2048 + C0) % 10000;
    boolean flag = true;
    string label = "rec2048";
};
type R2049 record {
    int id;
    int a = C1 + 2049;
    int b = C4 * 6;
    int c = (2049 + C1) % 10000;
    boolean flag = false;
    string label = "rec2049";
};
type R2050 record {
    int id;
    int a = C2 + 2050;
    int b = C7 * 7;
    int c = (2050 + C2) % 10000;
    boolean flag = true;
    string label = "rec2050";
};
type R2051 record {
    int id;
    int a = C3 + 2051;
    int b = C10 * 1;
    int c = (2051 + C3) % 10000;
    boolean flag = false;
    string label = "rec2051";
};
type R2052 record {
    int id;
    int a = C4 + 2052;
    int b = C13 * 2;
    int c = (2052 + C4) % 10000;
    boolean flag = true;
    string label = "rec2052";
};
type R2053 record {
    int id;
    int a = C5 + 2053;
    int b = C16 * 3;
    int c = (2053 + C5) % 10000;
    boolean flag = false;
    string label = "rec2053";
};
type R2054 record {
    int id;
    int a = C6 + 2054;
    int b = C19 * 4;
    int c = (2054 + C6) % 10000;
    boolean flag = true;
    string label = "rec2054";
};
type R2055 record {
    int id;
    int a = C7 + 2055;
    int b = C22 * 5;
    int c = (2055 + C7) % 10000;
    boolean flag = false;
    string label = "rec2055";
};
type R2056 record {
    int id;
    int a = C8 + 2056;
    int b = C25 * 6;
    int c = (2056 + C8) % 10000;
    boolean flag = true;
    string label = "rec2056";
};
type R2057 record {
    int id;
    int a = C9 + 2057;
    int b = C28 * 7;
    int c = (2057 + C9) % 10000;
    boolean flag = false;
    string label = "rec2057";
};
type R2058 record {
    int id;
    int a = C10 + 2058;
    int b = C31 * 1;
    int c = (2058 + C10) % 10000;
    boolean flag = true;
    string label = "rec2058";
};
type R2059 record {
    int id;
    int a = C11 + 2059;
    int b = C34 * 2;
    int c = (2059 + C11) % 10000;
    boolean flag = false;
    string label = "rec2059";
};
type R2060 record {
    int id;
    int a = C12 + 2060;
    int b = C37 * 3;
    int c = (2060 + C12) % 10000;
    boolean flag = true;
    string label = "rec2060";
};
type R2061 record {
    int id;
    int a = C13 + 2061;
    int b = C40 * 4;
    int c = (2061 + C13) % 10000;
    boolean flag = false;
    string label = "rec2061";
};
type R2062 record {
    int id;
    int a = C14 + 2062;
    int b = C43 * 5;
    int c = (2062 + C14) % 10000;
    boolean flag = true;
    string label = "rec2062";
};
type R2063 record {
    int id;
    int a = C15 + 2063;
    int b = C46 * 6;
    int c = (2063 + C15) % 10000;
    boolean flag = false;
    string label = "rec2063";
};
type R2064 record {
    int id;
    int a = C16 + 2064;
    int b = C49 * 7;
    int c = (2064 + C16) % 10000;
    boolean flag = true;
    string label = "rec2064";
};
type R2065 record {
    int id;
    int a = C17 + 2065;
    int b = C52 * 1;
    int c = (2065 + C17) % 10000;
    boolean flag = false;
    string label = "rec2065";
};
type R2066 record {
    int id;
    int a = C18 + 2066;
    int b = C55 * 2;
    int c = (2066 + C18) % 10000;
    boolean flag = true;
    string label = "rec2066";
};
type R2067 record {
    int id;
    int a = C19 + 2067;
    int b = C58 * 3;
    int c = (2067 + C19) % 10000;
    boolean flag = false;
    string label = "rec2067";
};
type R2068 record {
    int id;
    int a = C20 + 2068;
    int b = C61 * 4;
    int c = (2068 + C20) % 10000;
    boolean flag = true;
    string label = "rec2068";
};
type R2069 record {
    int id;
    int a = C21 + 2069;
    int b = C0 * 5;
    int c = (2069 + C21) % 10000;
    boolean flag = false;
    string label = "rec2069";
};
type R2070 record {
    int id;
    int a = C22 + 2070;
    int b = C3 * 6;
    int c = (2070 + C22) % 10000;
    boolean flag = true;
    string label = "rec2070";
};
type R2071 record {
    int id;
    int a = C23 + 2071;
    int b = C6 * 7;
    int c = (2071 + C23) % 10000;
    boolean flag = false;
    string label = "rec2071";
};
type R2072 record {
    int id;
    int a = C24 + 2072;
    int b = C9 * 1;
    int c = (2072 + C24) % 10000;
    boolean flag = true;
    string label = "rec2072";
};
type R2073 record {
    int id;
    int a = C25 + 2073;
    int b = C12 * 2;
    int c = (2073 + C25) % 10000;
    boolean flag = false;
    string label = "rec2073";
};
type R2074 record {
    int id;
    int a = C26 + 2074;
    int b = C15 * 3;
    int c = (2074 + C26) % 10000;
    boolean flag = true;
    string label = "rec2074";
};
type R2075 record {
    int id;
    int a = C27 + 2075;
    int b = C18 * 4;
    int c = (2075 + C27) % 10000;
    boolean flag = false;
    string label = "rec2075";
};
type R2076 record {
    int id;
    int a = C28 + 2076;
    int b = C21 * 5;
    int c = (2076 + C28) % 10000;
    boolean flag = true;
    string label = "rec2076";
};
type R2077 record {
    int id;
    int a = C29 + 2077;
    int b = C24 * 6;
    int c = (2077 + C29) % 10000;
    boolean flag = false;
    string label = "rec2077";
};
type R2078 record {
    int id;
    int a = C30 + 2078;
    int b = C27 * 7;
    int c = (2078 + C30) % 10000;
    boolean flag = true;
    string label = "rec2078";
};
type R2079 record {
    int id;
    int a = C31 + 2079;
    int b = C30 * 1;
    int c = (2079 + C31) % 10000;
    boolean flag = false;
    string label = "rec2079";
};
type R2080 record {
    int id;
    int a = C32 + 2080;
    int b = C33 * 2;
    int c = (2080 + C32) % 10000;
    boolean flag = true;
    string label = "rec2080";
};
type R2081 record {
    int id;
    int a = C33 + 2081;
    int b = C36 * 3;
    int c = (2081 + C33) % 10000;
    boolean flag = false;
    string label = "rec2081";
};
type R2082 record {
    int id;
    int a = C34 + 2082;
    int b = C39 * 4;
    int c = (2082 + C34) % 10000;
    boolean flag = true;
    string label = "rec2082";
};
type R2083 record {
    int id;
    int a = C35 + 2083;
    int b = C42 * 5;
    int c = (2083 + C35) % 10000;
    boolean flag = false;
    string label = "rec2083";
};
type R2084 record {
    int id;
    int a = C36 + 2084;
    int b = C45 * 6;
    int c = (2084 + C36) % 10000;
    boolean flag = true;
    string label = "rec2084";
};
type R2085 record {
    int id;
    int a = C37 + 2085;
    int b = C48 * 7;
    int c = (2085 + C37) % 10000;
    boolean flag = false;
    string label = "rec2085";
};
type R2086 record {
    int id;
    int a = C38 + 2086;
    int b = C51 * 1;
    int c = (2086 + C38) % 10000;
    boolean flag = true;
    string label = "rec2086";
};
type R2087 record {
    int id;
    int a = C39 + 2087;
    int b = C54 * 2;
    int c = (2087 + C39) % 10000;
    boolean flag = false;
    string label = "rec2087";
};
type R2088 record {
    int id;
    int a = C40 + 2088;
    int b = C57 * 3;
    int c = (2088 + C40) % 10000;
    boolean flag = true;
    string label = "rec2088";
};
type R2089 record {
    int id;
    int a = C41 + 2089;
    int b = C60 * 4;
    int c = (2089 + C41) % 10000;
    boolean flag = false;
    string label = "rec2089";
};
type R2090 record {
    int id;
    int a = C42 + 2090;
    int b = C63 * 5;
    int c = (2090 + C42) % 10000;
    boolean flag = true;
    string label = "rec2090";
};
type R2091 record {
    int id;
    int a = C43 + 2091;
    int b = C2 * 6;
    int c = (2091 + C43) % 10000;
    boolean flag = false;
    string label = "rec2091";
};
type R2092 record {
    int id;
    int a = C44 + 2092;
    int b = C5 * 7;
    int c = (2092 + C44) % 10000;
    boolean flag = true;
    string label = "rec2092";
};
type R2093 record {
    int id;
    int a = C45 + 2093;
    int b = C8 * 1;
    int c = (2093 + C45) % 10000;
    boolean flag = false;
    string label = "rec2093";
};
type R2094 record {
    int id;
    int a = C46 + 2094;
    int b = C11 * 2;
    int c = (2094 + C46) % 10000;
    boolean flag = true;
    string label = "rec2094";
};
type R2095 record {
    int id;
    int a = C47 + 2095;
    int b = C14 * 3;
    int c = (2095 + C47) % 10000;
    boolean flag = false;
    string label = "rec2095";
};
type R2096 record {
    int id;
    int a = C48 + 2096;
    int b = C17 * 4;
    int c = (2096 + C48) % 10000;
    boolean flag = true;
    string label = "rec2096";
};
type R2097 record {
    int id;
    int a = C49 + 2097;
    int b = C20 * 5;
    int c = (2097 + C49) % 10000;
    boolean flag = false;
    string label = "rec2097";
};
type R2098 record {
    int id;
    int a = C50 + 2098;
    int b = C23 * 6;
    int c = (2098 + C50) % 10000;
    boolean flag = true;
    string label = "rec2098";
};
type R2099 record {
    int id;
    int a = C51 + 2099;
    int b = C26 * 7;
    int c = (2099 + C51) % 10000;
    boolean flag = false;
    string label = "rec2099";
};
type R2100 record {
    int id;
    int a = C52 + 2100;
    int b = C29 * 1;
    int c = (2100 + C52) % 10000;
    boolean flag = true;
    string label = "rec2100";
};
type R2101 record {
    int id;
    int a = C53 + 2101;
    int b = C32 * 2;
    int c = (2101 + C53) % 10000;
    boolean flag = false;
    string label = "rec2101";
};
type R2102 record {
    int id;
    int a = C54 + 2102;
    int b = C35 * 3;
    int c = (2102 + C54) % 10000;
    boolean flag = true;
    string label = "rec2102";
};
type R2103 record {
    int id;
    int a = C55 + 2103;
    int b = C38 * 4;
    int c = (2103 + C55) % 10000;
    boolean flag = false;
    string label = "rec2103";
};
type R2104 record {
    int id;
    int a = C56 + 2104;
    int b = C41 * 5;
    int c = (2104 + C56) % 10000;
    boolean flag = true;
    string label = "rec2104";
};
type R2105 record {
    int id;
    int a = C57 + 2105;
    int b = C44 * 6;
    int c = (2105 + C57) % 10000;
    boolean flag = false;
    string label = "rec2105";
};
type R2106 record {
    int id;
    int a = C58 + 2106;
    int b = C47 * 7;
    int c = (2106 + C58) % 10000;
    boolean flag = true;
    string label = "rec2106";
};
type R2107 record {
    int id;
    int a = C59 + 2107;
    int b = C50 * 1;
    int c = (2107 + C59) % 10000;
    boolean flag = false;
    string label = "rec2107";
};
type R2108 record {
    int id;
    int a = C60 + 2108;
    int b = C53 * 2;
    int c = (2108 + C60) % 10000;
    boolean flag = true;
    string label = "rec2108";
};
type R2109 record {
    int id;
    int a = C61 + 2109;
    int b = C56 * 3;
    int c = (2109 + C61) % 10000;
    boolean flag = false;
    string label = "rec2109";
};
type R2110 record {
    int id;
    int a = C62 + 2110;
    int b = C59 * 4;
    int c = (2110 + C62) % 10000;
    boolean flag = true;
    string label = "rec2110";
};
type R2111 record {
    int id;
    int a = C63 + 2111;
    int b = C62 * 5;
    int c = (2111 + C63) % 10000;
    boolean flag = false;
    string label = "rec2111";
};
type R2112 record {
    int id;
    int a = C0 + 2112;
    int b = C1 * 6;
    int c = (2112 + C0) % 10000;
    boolean flag = true;
    string label = "rec2112";
};
type R2113 record {
    int id;
    int a = C1 + 2113;
    int b = C4 * 7;
    int c = (2113 + C1) % 10000;
    boolean flag = false;
    string label = "rec2113";
};
type R2114 record {
    int id;
    int a = C2 + 2114;
    int b = C7 * 1;
    int c = (2114 + C2) % 10000;
    boolean flag = true;
    string label = "rec2114";
};
type R2115 record {
    int id;
    int a = C3 + 2115;
    int b = C10 * 2;
    int c = (2115 + C3) % 10000;
    boolean flag = false;
    string label = "rec2115";
};
type R2116 record {
    int id;
    int a = C4 + 2116;
    int b = C13 * 3;
    int c = (2116 + C4) % 10000;
    boolean flag = true;
    string label = "rec2116";
};
type R2117 record {
    int id;
    int a = C5 + 2117;
    int b = C16 * 4;
    int c = (2117 + C5) % 10000;
    boolean flag = false;
    string label = "rec2117";
};
type R2118 record {
    int id;
    int a = C6 + 2118;
    int b = C19 * 5;
    int c = (2118 + C6) % 10000;
    boolean flag = true;
    string label = "rec2118";
};
type R2119 record {
    int id;
    int a = C7 + 2119;
    int b = C22 * 6;
    int c = (2119 + C7) % 10000;
    boolean flag = false;
    string label = "rec2119";
};
type R2120 record {
    int id;
    int a = C8 + 2120;
    int b = C25 * 7;
    int c = (2120 + C8) % 10000;
    boolean flag = true;
    string label = "rec2120";
};
type R2121 record {
    int id;
    int a = C9 + 2121;
    int b = C28 * 1;
    int c = (2121 + C9) % 10000;
    boolean flag = false;
    string label = "rec2121";
};
type R2122 record {
    int id;
    int a = C10 + 2122;
    int b = C31 * 2;
    int c = (2122 + C10) % 10000;
    boolean flag = true;
    string label = "rec2122";
};
type R2123 record {
    int id;
    int a = C11 + 2123;
    int b = C34 * 3;
    int c = (2123 + C11) % 10000;
    boolean flag = false;
    string label = "rec2123";
};
type R2124 record {
    int id;
    int a = C12 + 2124;
    int b = C37 * 4;
    int c = (2124 + C12) % 10000;
    boolean flag = true;
    string label = "rec2124";
};
type R2125 record {
    int id;
    int a = C13 + 2125;
    int b = C40 * 5;
    int c = (2125 + C13) % 10000;
    boolean flag = false;
    string label = "rec2125";
};
type R2126 record {
    int id;
    int a = C14 + 2126;
    int b = C43 * 6;
    int c = (2126 + C14) % 10000;
    boolean flag = true;
    string label = "rec2126";
};
type R2127 record {
    int id;
    int a = C15 + 2127;
    int b = C46 * 7;
    int c = (2127 + C15) % 10000;
    boolean flag = false;
    string label = "rec2127";
};
type R2128 record {
    int id;
    int a = C16 + 2128;
    int b = C49 * 1;
    int c = (2128 + C16) % 10000;
    boolean flag = true;
    string label = "rec2128";
};
type R2129 record {
    int id;
    int a = C17 + 2129;
    int b = C52 * 2;
    int c = (2129 + C17) % 10000;
    boolean flag = false;
    string label = "rec2129";
};
type R2130 record {
    int id;
    int a = C18 + 2130;
    int b = C55 * 3;
    int c = (2130 + C18) % 10000;
    boolean flag = true;
    string label = "rec2130";
};
type R2131 record {
    int id;
    int a = C19 + 2131;
    int b = C58 * 4;
    int c = (2131 + C19) % 10000;
    boolean flag = false;
    string label = "rec2131";
};
type R2132 record {
    int id;
    int a = C20 + 2132;
    int b = C61 * 5;
    int c = (2132 + C20) % 10000;
    boolean flag = true;
    string label = "rec2132";
};
type R2133 record {
    int id;
    int a = C21 + 2133;
    int b = C0 * 6;
    int c = (2133 + C21) % 10000;
    boolean flag = false;
    string label = "rec2133";
};
type R2134 record {
    int id;
    int a = C22 + 2134;
    int b = C3 * 7;
    int c = (2134 + C22) % 10000;
    boolean flag = true;
    string label = "rec2134";
};
type R2135 record {
    int id;
    int a = C23 + 2135;
    int b = C6 * 1;
    int c = (2135 + C23) % 10000;
    boolean flag = false;
    string label = "rec2135";
};
type R2136 record {
    int id;
    int a = C24 + 2136;
    int b = C9 * 2;
    int c = (2136 + C24) % 10000;
    boolean flag = true;
    string label = "rec2136";
};
type R2137 record {
    int id;
    int a = C25 + 2137;
    int b = C12 * 3;
    int c = (2137 + C25) % 10000;
    boolean flag = false;
    string label = "rec2137";
};
type R2138 record {
    int id;
    int a = C26 + 2138;
    int b = C15 * 4;
    int c = (2138 + C26) % 10000;
    boolean flag = true;
    string label = "rec2138";
};
type R2139 record {
    int id;
    int a = C27 + 2139;
    int b = C18 * 5;
    int c = (2139 + C27) % 10000;
    boolean flag = false;
    string label = "rec2139";
};
type R2140 record {
    int id;
    int a = C28 + 2140;
    int b = C21 * 6;
    int c = (2140 + C28) % 10000;
    boolean flag = true;
    string label = "rec2140";
};
type R2141 record {
    int id;
    int a = C29 + 2141;
    int b = C24 * 7;
    int c = (2141 + C29) % 10000;
    boolean flag = false;
    string label = "rec2141";
};
type R2142 record {
    int id;
    int a = C30 + 2142;
    int b = C27 * 1;
    int c = (2142 + C30) % 10000;
    boolean flag = true;
    string label = "rec2142";
};
type R2143 record {
    int id;
    int a = C31 + 2143;
    int b = C30 * 2;
    int c = (2143 + C31) % 10000;
    boolean flag = false;
    string label = "rec2143";
};
type R2144 record {
    int id;
    int a = C32 + 2144;
    int b = C33 * 3;
    int c = (2144 + C32) % 10000;
    boolean flag = true;
    string label = "rec2144";
};
type R2145 record {
    int id;
    int a = C33 + 2145;
    int b = C36 * 4;
    int c = (2145 + C33) % 10000;
    boolean flag = false;
    string label = "rec2145";
};
type R2146 record {
    int id;
    int a = C34 + 2146;
    int b = C39 * 5;
    int c = (2146 + C34) % 10000;
    boolean flag = true;
    string label = "rec2146";
};
type R2147 record {
    int id;
    int a = C35 + 2147;
    int b = C42 * 6;
    int c = (2147 + C35) % 10000;
    boolean flag = false;
    string label = "rec2147";
};
type R2148 record {
    int id;
    int a = C36 + 2148;
    int b = C45 * 7;
    int c = (2148 + C36) % 10000;
    boolean flag = true;
    string label = "rec2148";
};
type R2149 record {
    int id;
    int a = C37 + 2149;
    int b = C48 * 1;
    int c = (2149 + C37) % 10000;
    boolean flag = false;
    string label = "rec2149";
};
type R2150 record {
    int id;
    int a = C38 + 2150;
    int b = C51 * 2;
    int c = (2150 + C38) % 10000;
    boolean flag = true;
    string label = "rec2150";
};
type R2151 record {
    int id;
    int a = C39 + 2151;
    int b = C54 * 3;
    int c = (2151 + C39) % 10000;
    boolean flag = false;
    string label = "rec2151";
};
type R2152 record {
    int id;
    int a = C40 + 2152;
    int b = C57 * 4;
    int c = (2152 + C40) % 10000;
    boolean flag = true;
    string label = "rec2152";
};
type R2153 record {
    int id;
    int a = C41 + 2153;
    int b = C60 * 5;
    int c = (2153 + C41) % 10000;
    boolean flag = false;
    string label = "rec2153";
};
type R2154 record {
    int id;
    int a = C42 + 2154;
    int b = C63 * 6;
    int c = (2154 + C42) % 10000;
    boolean flag = true;
    string label = "rec2154";
};
type R2155 record {
    int id;
    int a = C43 + 2155;
    int b = C2 * 7;
    int c = (2155 + C43) % 10000;
    boolean flag = false;
    string label = "rec2155";
};
type R2156 record {
    int id;
    int a = C44 + 2156;
    int b = C5 * 1;
    int c = (2156 + C44) % 10000;
    boolean flag = true;
    string label = "rec2156";
};
type R2157 record {
    int id;
    int a = C45 + 2157;
    int b = C8 * 2;
    int c = (2157 + C45) % 10000;
    boolean flag = false;
    string label = "rec2157";
};
type R2158 record {
    int id;
    int a = C46 + 2158;
    int b = C11 * 3;
    int c = (2158 + C46) % 10000;
    boolean flag = true;
    string label = "rec2158";
};
type R2159 record {
    int id;
    int a = C47 + 2159;
    int b = C14 * 4;
    int c = (2159 + C47) % 10000;
    boolean flag = false;
    string label = "rec2159";
};
type R2160 record {
    int id;
    int a = C48 + 2160;
    int b = C17 * 5;
    int c = (2160 + C48) % 10000;
    boolean flag = true;
    string label = "rec2160";
};
type R2161 record {
    int id;
    int a = C49 + 2161;
    int b = C20 * 6;
    int c = (2161 + C49) % 10000;
    boolean flag = false;
    string label = "rec2161";
};
type R2162 record {
    int id;
    int a = C50 + 2162;
    int b = C23 * 7;
    int c = (2162 + C50) % 10000;
    boolean flag = true;
    string label = "rec2162";
};
type R2163 record {
    int id;
    int a = C51 + 2163;
    int b = C26 * 1;
    int c = (2163 + C51) % 10000;
    boolean flag = false;
    string label = "rec2163";
};
type R2164 record {
    int id;
    int a = C52 + 2164;
    int b = C29 * 2;
    int c = (2164 + C52) % 10000;
    boolean flag = true;
    string label = "rec2164";
};
type R2165 record {
    int id;
    int a = C53 + 2165;
    int b = C32 * 3;
    int c = (2165 + C53) % 10000;
    boolean flag = false;
    string label = "rec2165";
};
type R2166 record {
    int id;
    int a = C54 + 2166;
    int b = C35 * 4;
    int c = (2166 + C54) % 10000;
    boolean flag = true;
    string label = "rec2166";
};
type R2167 record {
    int id;
    int a = C55 + 2167;
    int b = C38 * 5;
    int c = (2167 + C55) % 10000;
    boolean flag = false;
    string label = "rec2167";
};
type R2168 record {
    int id;
    int a = C56 + 2168;
    int b = C41 * 6;
    int c = (2168 + C56) % 10000;
    boolean flag = true;
    string label = "rec2168";
};
type R2169 record {
    int id;
    int a = C57 + 2169;
    int b = C44 * 7;
    int c = (2169 + C57) % 10000;
    boolean flag = false;
    string label = "rec2169";
};
type R2170 record {
    int id;
    int a = C58 + 2170;
    int b = C47 * 1;
    int c = (2170 + C58) % 10000;
    boolean flag = true;
    string label = "rec2170";
};
type R2171 record {
    int id;
    int a = C59 + 2171;
    int b = C50 * 2;
    int c = (2171 + C59) % 10000;
    boolean flag = false;
    string label = "rec2171";
};
type R2172 record {
    int id;
    int a = C60 + 2172;
    int b = C53 * 3;
    int c = (2172 + C60) % 10000;
    boolean flag = true;
    string label = "rec2172";
};
type R2173 record {
    int id;
    int a = C61 + 2173;
    int b = C56 * 4;
    int c = (2173 + C61) % 10000;
    boolean flag = false;
    string label = "rec2173";
};
type R2174 record {
    int id;
    int a = C62 + 2174;
    int b = C59 * 5;
    int c = (2174 + C62) % 10000;
    boolean flag = true;
    string label = "rec2174";
};
type R2175 record {
    int id;
    int a = C63 + 2175;
    int b = C62 * 6;
    int c = (2175 + C63) % 10000;
    boolean flag = false;
    string label = "rec2175";
};
type R2176 record {
    int id;
    int a = C0 + 2176;
    int b = C1 * 7;
    int c = (2176 + C0) % 10000;
    boolean flag = true;
    string label = "rec2176";
};
type R2177 record {
    int id;
    int a = C1 + 2177;
    int b = C4 * 1;
    int c = (2177 + C1) % 10000;
    boolean flag = false;
    string label = "rec2177";
};
type R2178 record {
    int id;
    int a = C2 + 2178;
    int b = C7 * 2;
    int c = (2178 + C2) % 10000;
    boolean flag = true;
    string label = "rec2178";
};
type R2179 record {
    int id;
    int a = C3 + 2179;
    int b = C10 * 3;
    int c = (2179 + C3) % 10000;
    boolean flag = false;
    string label = "rec2179";
};
type R2180 record {
    int id;
    int a = C4 + 2180;
    int b = C13 * 4;
    int c = (2180 + C4) % 10000;
    boolean flag = true;
    string label = "rec2180";
};
type R2181 record {
    int id;
    int a = C5 + 2181;
    int b = C16 * 5;
    int c = (2181 + C5) % 10000;
    boolean flag = false;
    string label = "rec2181";
};
type R2182 record {
    int id;
    int a = C6 + 2182;
    int b = C19 * 6;
    int c = (2182 + C6) % 10000;
    boolean flag = true;
    string label = "rec2182";
};
type R2183 record {
    int id;
    int a = C7 + 2183;
    int b = C22 * 7;
    int c = (2183 + C7) % 10000;
    boolean flag = false;
    string label = "rec2183";
};
type R2184 record {
    int id;
    int a = C8 + 2184;
    int b = C25 * 1;
    int c = (2184 + C8) % 10000;
    boolean flag = true;
    string label = "rec2184";
};
type R2185 record {
    int id;
    int a = C9 + 2185;
    int b = C28 * 2;
    int c = (2185 + C9) % 10000;
    boolean flag = false;
    string label = "rec2185";
};
type R2186 record {
    int id;
    int a = C10 + 2186;
    int b = C31 * 3;
    int c = (2186 + C10) % 10000;
    boolean flag = true;
    string label = "rec2186";
};
type R2187 record {
    int id;
    int a = C11 + 2187;
    int b = C34 * 4;
    int c = (2187 + C11) % 10000;
    boolean flag = false;
    string label = "rec2187";
};
type R2188 record {
    int id;
    int a = C12 + 2188;
    int b = C37 * 5;
    int c = (2188 + C12) % 10000;
    boolean flag = true;
    string label = "rec2188";
};
type R2189 record {
    int id;
    int a = C13 + 2189;
    int b = C40 * 6;
    int c = (2189 + C13) % 10000;
    boolean flag = false;
    string label = "rec2189";
};
type R2190 record {
    int id;
    int a = C14 + 2190;
    int b = C43 * 7;
    int c = (2190 + C14) % 10000;
    boolean flag = true;
    string label = "rec2190";
};
type R2191 record {
    int id;
    int a = C15 + 2191;
    int b = C46 * 1;
    int c = (2191 + C15) % 10000;
    boolean flag = false;
    string label = "rec2191";
};
type R2192 record {
    int id;
    int a = C16 + 2192;
    int b = C49 * 2;
    int c = (2192 + C16) % 10000;
    boolean flag = true;
    string label = "rec2192";
};
type R2193 record {
    int id;
    int a = C17 + 2193;
    int b = C52 * 3;
    int c = (2193 + C17) % 10000;
    boolean flag = false;
    string label = "rec2193";
};
type R2194 record {
    int id;
    int a = C18 + 2194;
    int b = C55 * 4;
    int c = (2194 + C18) % 10000;
    boolean flag = true;
    string label = "rec2194";
};
type R2195 record {
    int id;
    int a = C19 + 2195;
    int b = C58 * 5;
    int c = (2195 + C19) % 10000;
    boolean flag = false;
    string label = "rec2195";
};
type R2196 record {
    int id;
    int a = C20 + 2196;
    int b = C61 * 6;
    int c = (2196 + C20) % 10000;
    boolean flag = true;
    string label = "rec2196";
};
type R2197 record {
    int id;
    int a = C21 + 2197;
    int b = C0 * 7;
    int c = (2197 + C21) % 10000;
    boolean flag = false;
    string label = "rec2197";
};
type R2198 record {
    int id;
    int a = C22 + 2198;
    int b = C3 * 1;
    int c = (2198 + C22) % 10000;
    boolean flag = true;
    string label = "rec2198";
};
type R2199 record {
    int id;
    int a = C23 + 2199;
    int b = C6 * 2;
    int c = (2199 + C23) % 10000;
    boolean flag = false;
    string label = "rec2199";
};
type R2200 record {
    int id;
    int a = C24 + 2200;
    int b = C9 * 3;
    int c = (2200 + C24) % 10000;
    boolean flag = true;
    string label = "rec2200";
};
type R2201 record {
    int id;
    int a = C25 + 2201;
    int b = C12 * 4;
    int c = (2201 + C25) % 10000;
    boolean flag = false;
    string label = "rec2201";
};
type R2202 record {
    int id;
    int a = C26 + 2202;
    int b = C15 * 5;
    int c = (2202 + C26) % 10000;
    boolean flag = true;
    string label = "rec2202";
};
type R2203 record {
    int id;
    int a = C27 + 2203;
    int b = C18 * 6;
    int c = (2203 + C27) % 10000;
    boolean flag = false;
    string label = "rec2203";
};
type R2204 record {
    int id;
    int a = C28 + 2204;
    int b = C21 * 7;
    int c = (2204 + C28) % 10000;
    boolean flag = true;
    string label = "rec2204";
};
type R2205 record {
    int id;
    int a = C29 + 2205;
    int b = C24 * 1;
    int c = (2205 + C29) % 10000;
    boolean flag = false;
    string label = "rec2205";
};
type R2206 record {
    int id;
    int a = C30 + 2206;
    int b = C27 * 2;
    int c = (2206 + C30) % 10000;
    boolean flag = true;
    string label = "rec2206";
};
type R2207 record {
    int id;
    int a = C31 + 2207;
    int b = C30 * 3;
    int c = (2207 + C31) % 10000;
    boolean flag = false;
    string label = "rec2207";
};
type R2208 record {
    int id;
    int a = C32 + 2208;
    int b = C33 * 4;
    int c = (2208 + C32) % 10000;
    boolean flag = true;
    string label = "rec2208";
};
type R2209 record {
    int id;
    int a = C33 + 2209;
    int b = C36 * 5;
    int c = (2209 + C33) % 10000;
    boolean flag = false;
    string label = "rec2209";
};
type R2210 record {
    int id;
    int a = C34 + 2210;
    int b = C39 * 6;
    int c = (2210 + C34) % 10000;
    boolean flag = true;
    string label = "rec2210";
};
type R2211 record {
    int id;
    int a = C35 + 2211;
    int b = C42 * 7;
    int c = (2211 + C35) % 10000;
    boolean flag = false;
    string label = "rec2211";
};
type R2212 record {
    int id;
    int a = C36 + 2212;
    int b = C45 * 1;
    int c = (2212 + C36) % 10000;
    boolean flag = true;
    string label = "rec2212";
};
type R2213 record {
    int id;
    int a = C37 + 2213;
    int b = C48 * 2;
    int c = (2213 + C37) % 10000;
    boolean flag = false;
    string label = "rec2213";
};
type R2214 record {
    int id;
    int a = C38 + 2214;
    int b = C51 * 3;
    int c = (2214 + C38) % 10000;
    boolean flag = true;
    string label = "rec2214";
};
type R2215 record {
    int id;
    int a = C39 + 2215;
    int b = C54 * 4;
    int c = (2215 + C39) % 10000;
    boolean flag = false;
    string label = "rec2215";
};
type R2216 record {
    int id;
    int a = C40 + 2216;
    int b = C57 * 5;
    int c = (2216 + C40) % 10000;
    boolean flag = true;
    string label = "rec2216";
};
type R2217 record {
    int id;
    int a = C41 + 2217;
    int b = C60 * 6;
    int c = (2217 + C41) % 10000;
    boolean flag = false;
    string label = "rec2217";
};
type R2218 record {
    int id;
    int a = C42 + 2218;
    int b = C63 * 7;
    int c = (2218 + C42) % 10000;
    boolean flag = true;
    string label = "rec2218";
};
type R2219 record {
    int id;
    int a = C43 + 2219;
    int b = C2 * 1;
    int c = (2219 + C43) % 10000;
    boolean flag = false;
    string label = "rec2219";
};
type R2220 record {
    int id;
    int a = C44 + 2220;
    int b = C5 * 2;
    int c = (2220 + C44) % 10000;
    boolean flag = true;
    string label = "rec2220";
};
type R2221 record {
    int id;
    int a = C45 + 2221;
    int b = C8 * 3;
    int c = (2221 + C45) % 10000;
    boolean flag = false;
    string label = "rec2221";
};
type R2222 record {
    int id;
    int a = C46 + 2222;
    int b = C11 * 4;
    int c = (2222 + C46) % 10000;
    boolean flag = true;
    string label = "rec2222";
};
type R2223 record {
    int id;
    int a = C47 + 2223;
    int b = C14 * 5;
    int c = (2223 + C47) % 10000;
    boolean flag = false;
    string label = "rec2223";
};
type R2224 record {
    int id;
    int a = C48 + 2224;
    int b = C17 * 6;
    int c = (2224 + C48) % 10000;
    boolean flag = true;
    string label = "rec2224";
};
type R2225 record {
    int id;
    int a = C49 + 2225;
    int b = C20 * 7;
    int c = (2225 + C49) % 10000;
    boolean flag = false;
    string label = "rec2225";
};
type R2226 record {
    int id;
    int a = C50 + 2226;
    int b = C23 * 1;
    int c = (2226 + C50) % 10000;
    boolean flag = true;
    string label = "rec2226";
};
type R2227 record {
    int id;
    int a = C51 + 2227;
    int b = C26 * 2;
    int c = (2227 + C51) % 10000;
    boolean flag = false;
    string label = "rec2227";
};
type R2228 record {
    int id;
    int a = C52 + 2228;
    int b = C29 * 3;
    int c = (2228 + C52) % 10000;
    boolean flag = true;
    string label = "rec2228";
};
type R2229 record {
    int id;
    int a = C53 + 2229;
    int b = C32 * 4;
    int c = (2229 + C53) % 10000;
    boolean flag = false;
    string label = "rec2229";
};
type R2230 record {
    int id;
    int a = C54 + 2230;
    int b = C35 * 5;
    int c = (2230 + C54) % 10000;
    boolean flag = true;
    string label = "rec2230";
};
type R2231 record {
    int id;
    int a = C55 + 2231;
    int b = C38 * 6;
    int c = (2231 + C55) % 10000;
    boolean flag = false;
    string label = "rec2231";
};
type R2232 record {
    int id;
    int a = C56 + 2232;
    int b = C41 * 7;
    int c = (2232 + C56) % 10000;
    boolean flag = true;
    string label = "rec2232";
};
type R2233 record {
    int id;
    int a = C57 + 2233;
    int b = C44 * 1;
    int c = (2233 + C57) % 10000;
    boolean flag = false;
    string label = "rec2233";
};
type R2234 record {
    int id;
    int a = C58 + 2234;
    int b = C47 * 2;
    int c = (2234 + C58) % 10000;
    boolean flag = true;
    string label = "rec2234";
};
type R2235 record {
    int id;
    int a = C59 + 2235;
    int b = C50 * 3;
    int c = (2235 + C59) % 10000;
    boolean flag = false;
    string label = "rec2235";
};
type R2236 record {
    int id;
    int a = C60 + 2236;
    int b = C53 * 4;
    int c = (2236 + C60) % 10000;
    boolean flag = true;
    string label = "rec2236";
};
type R2237 record {
    int id;
    int a = C61 + 2237;
    int b = C56 * 5;
    int c = (2237 + C61) % 10000;
    boolean flag = false;
    string label = "rec2237";
};
type R2238 record {
    int id;
    int a = C62 + 2238;
    int b = C59 * 6;
    int c = (2238 + C62) % 10000;
    boolean flag = true;
    string label = "rec2238";
};
type R2239 record {
    int id;
    int a = C63 + 2239;
    int b = C62 * 7;
    int c = (2239 + C63) % 10000;
    boolean flag = false;
    string label = "rec2239";
};
type R2240 record {
    int id;
    int a = C0 + 2240;
    int b = C1 * 1;
    int c = (2240 + C0) % 10000;
    boolean flag = true;
    string label = "rec2240";
};
type R2241 record {
    int id;
    int a = C1 + 2241;
    int b = C4 * 2;
    int c = (2241 + C1) % 10000;
    boolean flag = false;
    string label = "rec2241";
};
type R2242 record {
    int id;
    int a = C2 + 2242;
    int b = C7 * 3;
    int c = (2242 + C2) % 10000;
    boolean flag = true;
    string label = "rec2242";
};
type R2243 record {
    int id;
    int a = C3 + 2243;
    int b = C10 * 4;
    int c = (2243 + C3) % 10000;
    boolean flag = false;
    string label = "rec2243";
};
type R2244 record {
    int id;
    int a = C4 + 2244;
    int b = C13 * 5;
    int c = (2244 + C4) % 10000;
    boolean flag = true;
    string label = "rec2244";
};
type R2245 record {
    int id;
    int a = C5 + 2245;
    int b = C16 * 6;
    int c = (2245 + C5) % 10000;
    boolean flag = false;
    string label = "rec2245";
};
type R2246 record {
    int id;
    int a = C6 + 2246;
    int b = C19 * 7;
    int c = (2246 + C6) % 10000;
    boolean flag = true;
    string label = "rec2246";
};
type R2247 record {
    int id;
    int a = C7 + 2247;
    int b = C22 * 1;
    int c = (2247 + C7) % 10000;
    boolean flag = false;
    string label = "rec2247";
};
type R2248 record {
    int id;
    int a = C8 + 2248;
    int b = C25 * 2;
    int c = (2248 + C8) % 10000;
    boolean flag = true;
    string label = "rec2248";
};
type R2249 record {
    int id;
    int a = C9 + 2249;
    int b = C28 * 3;
    int c = (2249 + C9) % 10000;
    boolean flag = false;
    string label = "rec2249";
};
type R2250 record {
    int id;
    int a = C10 + 2250;
    int b = C31 * 4;
    int c = (2250 + C10) % 10000;
    boolean flag = true;
    string label = "rec2250";
};
type R2251 record {
    int id;
    int a = C11 + 2251;
    int b = C34 * 5;
    int c = (2251 + C11) % 10000;
    boolean flag = false;
    string label = "rec2251";
};
type R2252 record {
    int id;
    int a = C12 + 2252;
    int b = C37 * 6;
    int c = (2252 + C12) % 10000;
    boolean flag = true;
    string label = "rec2252";
};
type R2253 record {
    int id;
    int a = C13 + 2253;
    int b = C40 * 7;
    int c = (2253 + C13) % 10000;
    boolean flag = false;
    string label = "rec2253";
};
type R2254 record {
    int id;
    int a = C14 + 2254;
    int b = C43 * 1;
    int c = (2254 + C14) % 10000;
    boolean flag = true;
    string label = "rec2254";
};
type R2255 record {
    int id;
    int a = C15 + 2255;
    int b = C46 * 2;
    int c = (2255 + C15) % 10000;
    boolean flag = false;
    string label = "rec2255";
};
type R2256 record {
    int id;
    int a = C16 + 2256;
    int b = C49 * 3;
    int c = (2256 + C16) % 10000;
    boolean flag = true;
    string label = "rec2256";
};
type R2257 record {
    int id;
    int a = C17 + 2257;
    int b = C52 * 4;
    int c = (2257 + C17) % 10000;
    boolean flag = false;
    string label = "rec2257";
};
type R2258 record {
    int id;
    int a = C18 + 2258;
    int b = C55 * 5;
    int c = (2258 + C18) % 10000;
    boolean flag = true;
    string label = "rec2258";
};
type R2259 record {
    int id;
    int a = C19 + 2259;
    int b = C58 * 6;
    int c = (2259 + C19) % 10000;
    boolean flag = false;
    string label = "rec2259";
};
type R2260 record {
    int id;
    int a = C20 + 2260;
    int b = C61 * 7;
    int c = (2260 + C20) % 10000;
    boolean flag = true;
    string label = "rec2260";
};
type R2261 record {
    int id;
    int a = C21 + 2261;
    int b = C0 * 1;
    int c = (2261 + C21) % 10000;
    boolean flag = false;
    string label = "rec2261";
};
type R2262 record {
    int id;
    int a = C22 + 2262;
    int b = C3 * 2;
    int c = (2262 + C22) % 10000;
    boolean flag = true;
    string label = "rec2262";
};
type R2263 record {
    int id;
    int a = C23 + 2263;
    int b = C6 * 3;
    int c = (2263 + C23) % 10000;
    boolean flag = false;
    string label = "rec2263";
};
type R2264 record {
    int id;
    int a = C24 + 2264;
    int b = C9 * 4;
    int c = (2264 + C24) % 10000;
    boolean flag = true;
    string label = "rec2264";
};
type R2265 record {
    int id;
    int a = C25 + 2265;
    int b = C12 * 5;
    int c = (2265 + C25) % 10000;
    boolean flag = false;
    string label = "rec2265";
};
type R2266 record {
    int id;
    int a = C26 + 2266;
    int b = C15 * 6;
    int c = (2266 + C26) % 10000;
    boolean flag = true;
    string label = "rec2266";
};
type R2267 record {
    int id;
    int a = C27 + 2267;
    int b = C18 * 7;
    int c = (2267 + C27) % 10000;
    boolean flag = false;
    string label = "rec2267";
};
type R2268 record {
    int id;
    int a = C28 + 2268;
    int b = C21 * 1;
    int c = (2268 + C28) % 10000;
    boolean flag = true;
    string label = "rec2268";
};
type R2269 record {
    int id;
    int a = C29 + 2269;
    int b = C24 * 2;
    int c = (2269 + C29) % 10000;
    boolean flag = false;
    string label = "rec2269";
};
type R2270 record {
    int id;
    int a = C30 + 2270;
    int b = C27 * 3;
    int c = (2270 + C30) % 10000;
    boolean flag = true;
    string label = "rec2270";
};
type R2271 record {
    int id;
    int a = C31 + 2271;
    int b = C30 * 4;
    int c = (2271 + C31) % 10000;
    boolean flag = false;
    string label = "rec2271";
};
type R2272 record {
    int id;
    int a = C32 + 2272;
    int b = C33 * 5;
    int c = (2272 + C32) % 10000;
    boolean flag = true;
    string label = "rec2272";
};
type R2273 record {
    int id;
    int a = C33 + 2273;
    int b = C36 * 6;
    int c = (2273 + C33) % 10000;
    boolean flag = false;
    string label = "rec2273";
};
type R2274 record {
    int id;
    int a = C34 + 2274;
    int b = C39 * 7;
    int c = (2274 + C34) % 10000;
    boolean flag = true;
    string label = "rec2274";
};
type R2275 record {
    int id;
    int a = C35 + 2275;
    int b = C42 * 1;
    int c = (2275 + C35) % 10000;
    boolean flag = false;
    string label = "rec2275";
};
type R2276 record {
    int id;
    int a = C36 + 2276;
    int b = C45 * 2;
    int c = (2276 + C36) % 10000;
    boolean flag = true;
    string label = "rec2276";
};
type R2277 record {
    int id;
    int a = C37 + 2277;
    int b = C48 * 3;
    int c = (2277 + C37) % 10000;
    boolean flag = false;
    string label = "rec2277";
};
type R2278 record {
    int id;
    int a = C38 + 2278;
    int b = C51 * 4;
    int c = (2278 + C38) % 10000;
    boolean flag = true;
    string label = "rec2278";
};
type R2279 record {
    int id;
    int a = C39 + 2279;
    int b = C54 * 5;
    int c = (2279 + C39) % 10000;
    boolean flag = false;
    string label = "rec2279";
};
type R2280 record {
    int id;
    int a = C40 + 2280;
    int b = C57 * 6;
    int c = (2280 + C40) % 10000;
    boolean flag = true;
    string label = "rec2280";
};
type R2281 record {
    int id;
    int a = C41 + 2281;
    int b = C60 * 7;
    int c = (2281 + C41) % 10000;
    boolean flag = false;
    string label = "rec2281";
};
type R2282 record {
    int id;
    int a = C42 + 2282;
    int b = C63 * 1;
    int c = (2282 + C42) % 10000;
    boolean flag = true;
    string label = "rec2282";
};
type R2283 record {
    int id;
    int a = C43 + 2283;
    int b = C2 * 2;
    int c = (2283 + C43) % 10000;
    boolean flag = false;
    string label = "rec2283";
};
type R2284 record {
    int id;
    int a = C44 + 2284;
    int b = C5 * 3;
    int c = (2284 + C44) % 10000;
    boolean flag = true;
    string label = "rec2284";
};
type R2285 record {
    int id;
    int a = C45 + 2285;
    int b = C8 * 4;
    int c = (2285 + C45) % 10000;
    boolean flag = false;
    string label = "rec2285";
};
type R2286 record {
    int id;
    int a = C46 + 2286;
    int b = C11 * 5;
    int c = (2286 + C46) % 10000;
    boolean flag = true;
    string label = "rec2286";
};
type R2287 record {
    int id;
    int a = C47 + 2287;
    int b = C14 * 6;
    int c = (2287 + C47) % 10000;
    boolean flag = false;
    string label = "rec2287";
};
type R2288 record {
    int id;
    int a = C48 + 2288;
    int b = C17 * 7;
    int c = (2288 + C48) % 10000;
    boolean flag = true;
    string label = "rec2288";
};
type R2289 record {
    int id;
    int a = C49 + 2289;
    int b = C20 * 1;
    int c = (2289 + C49) % 10000;
    boolean flag = false;
    string label = "rec2289";
};
type R2290 record {
    int id;
    int a = C50 + 2290;
    int b = C23 * 2;
    int c = (2290 + C50) % 10000;
    boolean flag = true;
    string label = "rec2290";
};
type R2291 record {
    int id;
    int a = C51 + 2291;
    int b = C26 * 3;
    int c = (2291 + C51) % 10000;
    boolean flag = false;
    string label = "rec2291";
};
type R2292 record {
    int id;
    int a = C52 + 2292;
    int b = C29 * 4;
    int c = (2292 + C52) % 10000;
    boolean flag = true;
    string label = "rec2292";
};
type R2293 record {
    int id;
    int a = C53 + 2293;
    int b = C32 * 5;
    int c = (2293 + C53) % 10000;
    boolean flag = false;
    string label = "rec2293";
};
type R2294 record {
    int id;
    int a = C54 + 2294;
    int b = C35 * 6;
    int c = (2294 + C54) % 10000;
    boolean flag = true;
    string label = "rec2294";
};
type R2295 record {
    int id;
    int a = C55 + 2295;
    int b = C38 * 7;
    int c = (2295 + C55) % 10000;
    boolean flag = false;
    string label = "rec2295";
};
type R2296 record {
    int id;
    int a = C56 + 2296;
    int b = C41 * 1;
    int c = (2296 + C56) % 10000;
    boolean flag = true;
    string label = "rec2296";
};
type R2297 record {
    int id;
    int a = C57 + 2297;
    int b = C44 * 2;
    int c = (2297 + C57) % 10000;
    boolean flag = false;
    string label = "rec2297";
};
type R2298 record {
    int id;
    int a = C58 + 2298;
    int b = C47 * 3;
    int c = (2298 + C58) % 10000;
    boolean flag = true;
    string label = "rec2298";
};
type R2299 record {
    int id;
    int a = C59 + 2299;
    int b = C50 * 4;
    int c = (2299 + C59) % 10000;
    boolean flag = false;
    string label = "rec2299";
};
type R2300 record {
    int id;
    int a = C60 + 2300;
    int b = C53 * 5;
    int c = (2300 + C60) % 10000;
    boolean flag = true;
    string label = "rec2300";
};
type R2301 record {
    int id;
    int a = C61 + 2301;
    int b = C56 * 6;
    int c = (2301 + C61) % 10000;
    boolean flag = false;
    string label = "rec2301";
};
type R2302 record {
    int id;
    int a = C62 + 2302;
    int b = C59 * 7;
    int c = (2302 + C62) % 10000;
    boolean flag = true;
    string label = "rec2302";
};
type R2303 record {
    int id;
    int a = C63 + 2303;
    int b = C62 * 1;
    int c = (2303 + C63) % 10000;
    boolean flag = false;
    string label = "rec2303";
};
type R2304 record {
    int id;
    int a = C0 + 2304;
    int b = C1 * 2;
    int c = (2304 + C0) % 10000;
    boolean flag = true;
    string label = "rec2304";
};
type R2305 record {
    int id;
    int a = C1 + 2305;
    int b = C4 * 3;
    int c = (2305 + C1) % 10000;
    boolean flag = false;
    string label = "rec2305";
};
type R2306 record {
    int id;
    int a = C2 + 2306;
    int b = C7 * 4;
    int c = (2306 + C2) % 10000;
    boolean flag = true;
    string label = "rec2306";
};
type R2307 record {
    int id;
    int a = C3 + 2307;
    int b = C10 * 5;
    int c = (2307 + C3) % 10000;
    boolean flag = false;
    string label = "rec2307";
};
type R2308 record {
    int id;
    int a = C4 + 2308;
    int b = C13 * 6;
    int c = (2308 + C4) % 10000;
    boolean flag = true;
    string label = "rec2308";
};
type R2309 record {
    int id;
    int a = C5 + 2309;
    int b = C16 * 7;
    int c = (2309 + C5) % 10000;
    boolean flag = false;
    string label = "rec2309";
};
type R2310 record {
    int id;
    int a = C6 + 2310;
    int b = C19 * 1;
    int c = (2310 + C6) % 10000;
    boolean flag = true;
    string label = "rec2310";
};
type R2311 record {
    int id;
    int a = C7 + 2311;
    int b = C22 * 2;
    int c = (2311 + C7) % 10000;
    boolean flag = false;
    string label = "rec2311";
};
type R2312 record {
    int id;
    int a = C8 + 2312;
    int b = C25 * 3;
    int c = (2312 + C8) % 10000;
    boolean flag = true;
    string label = "rec2312";
};
type R2313 record {
    int id;
    int a = C9 + 2313;
    int b = C28 * 4;
    int c = (2313 + C9) % 10000;
    boolean flag = false;
    string label = "rec2313";
};
type R2314 record {
    int id;
    int a = C10 + 2314;
    int b = C31 * 5;
    int c = (2314 + C10) % 10000;
    boolean flag = true;
    string label = "rec2314";
};
type R2315 record {
    int id;
    int a = C11 + 2315;
    int b = C34 * 6;
    int c = (2315 + C11) % 10000;
    boolean flag = false;
    string label = "rec2315";
};
type R2316 record {
    int id;
    int a = C12 + 2316;
    int b = C37 * 7;
    int c = (2316 + C12) % 10000;
    boolean flag = true;
    string label = "rec2316";
};
type R2317 record {
    int id;
    int a = C13 + 2317;
    int b = C40 * 1;
    int c = (2317 + C13) % 10000;
    boolean flag = false;
    string label = "rec2317";
};
type R2318 record {
    int id;
    int a = C14 + 2318;
    int b = C43 * 2;
    int c = (2318 + C14) % 10000;
    boolean flag = true;
    string label = "rec2318";
};
type R2319 record {
    int id;
    int a = C15 + 2319;
    int b = C46 * 3;
    int c = (2319 + C15) % 10000;
    boolean flag = false;
    string label = "rec2319";
};
type R2320 record {
    int id;
    int a = C16 + 2320;
    int b = C49 * 4;
    int c = (2320 + C16) % 10000;
    boolean flag = true;
    string label = "rec2320";
};
type R2321 record {
    int id;
    int a = C17 + 2321;
    int b = C52 * 5;
    int c = (2321 + C17) % 10000;
    boolean flag = false;
    string label = "rec2321";
};
type R2322 record {
    int id;
    int a = C18 + 2322;
    int b = C55 * 6;
    int c = (2322 + C18) % 10000;
    boolean flag = true;
    string label = "rec2322";
};
type R2323 record {
    int id;
    int a = C19 + 2323;
    int b = C58 * 7;
    int c = (2323 + C19) % 10000;
    boolean flag = false;
    string label = "rec2323";
};
type R2324 record {
    int id;
    int a = C20 + 2324;
    int b = C61 * 1;
    int c = (2324 + C20) % 10000;
    boolean flag = true;
    string label = "rec2324";
};
type R2325 record {
    int id;
    int a = C21 + 2325;
    int b = C0 * 2;
    int c = (2325 + C21) % 10000;
    boolean flag = false;
    string label = "rec2325";
};
type R2326 record {
    int id;
    int a = C22 + 2326;
    int b = C3 * 3;
    int c = (2326 + C22) % 10000;
    boolean flag = true;
    string label = "rec2326";
};
type R2327 record {
    int id;
    int a = C23 + 2327;
    int b = C6 * 4;
    int c = (2327 + C23) % 10000;
    boolean flag = false;
    string label = "rec2327";
};
type R2328 record {
    int id;
    int a = C24 + 2328;
    int b = C9 * 5;
    int c = (2328 + C24) % 10000;
    boolean flag = true;
    string label = "rec2328";
};
type R2329 record {
    int id;
    int a = C25 + 2329;
    int b = C12 * 6;
    int c = (2329 + C25) % 10000;
    boolean flag = false;
    string label = "rec2329";
};
type R2330 record {
    int id;
    int a = C26 + 2330;
    int b = C15 * 7;
    int c = (2330 + C26) % 10000;
    boolean flag = true;
    string label = "rec2330";
};
type R2331 record {
    int id;
    int a = C27 + 2331;
    int b = C18 * 1;
    int c = (2331 + C27) % 10000;
    boolean flag = false;
    string label = "rec2331";
};
type R2332 record {
    int id;
    int a = C28 + 2332;
    int b = C21 * 2;
    int c = (2332 + C28) % 10000;
    boolean flag = true;
    string label = "rec2332";
};
type R2333 record {
    int id;
    int a = C29 + 2333;
    int b = C24 * 3;
    int c = (2333 + C29) % 10000;
    boolean flag = false;
    string label = "rec2333";
};
type R2334 record {
    int id;
    int a = C30 + 2334;
    int b = C27 * 4;
    int c = (2334 + C30) % 10000;
    boolean flag = true;
    string label = "rec2334";
};
type R2335 record {
    int id;
    int a = C31 + 2335;
    int b = C30 * 5;
    int c = (2335 + C31) % 10000;
    boolean flag = false;
    string label = "rec2335";
};
type R2336 record {
    int id;
    int a = C32 + 2336;
    int b = C33 * 6;
    int c = (2336 + C32) % 10000;
    boolean flag = true;
    string label = "rec2336";
};
type R2337 record {
    int id;
    int a = C33 + 2337;
    int b = C36 * 7;
    int c = (2337 + C33) % 10000;
    boolean flag = false;
    string label = "rec2337";
};
type R2338 record {
    int id;
    int a = C34 + 2338;
    int b = C39 * 1;
    int c = (2338 + C34) % 10000;
    boolean flag = true;
    string label = "rec2338";
};
type R2339 record {
    int id;
    int a = C35 + 2339;
    int b = C42 * 2;
    int c = (2339 + C35) % 10000;
    boolean flag = false;
    string label = "rec2339";
};
type R2340 record {
    int id;
    int a = C36 + 2340;
    int b = C45 * 3;
    int c = (2340 + C36) % 10000;
    boolean flag = true;
    string label = "rec2340";
};
type R2341 record {
    int id;
    int a = C37 + 2341;
    int b = C48 * 4;
    int c = (2341 + C37) % 10000;
    boolean flag = false;
    string label = "rec2341";
};
type R2342 record {
    int id;
    int a = C38 + 2342;
    int b = C51 * 5;
    int c = (2342 + C38) % 10000;
    boolean flag = true;
    string label = "rec2342";
};
type R2343 record {
    int id;
    int a = C39 + 2343;
    int b = C54 * 6;
    int c = (2343 + C39) % 10000;
    boolean flag = false;
    string label = "rec2343";
};
type R2344 record {
    int id;
    int a = C40 + 2344;
    int b = C57 * 7;
    int c = (2344 + C40) % 10000;
    boolean flag = true;
    string label = "rec2344";
};
type R2345 record {
    int id;
    int a = C41 + 2345;
    int b = C60 * 1;
    int c = (2345 + C41) % 10000;
    boolean flag = false;
    string label = "rec2345";
};
type R2346 record {
    int id;
    int a = C42 + 2346;
    int b = C63 * 2;
    int c = (2346 + C42) % 10000;
    boolean flag = true;
    string label = "rec2346";
};
type R2347 record {
    int id;
    int a = C43 + 2347;
    int b = C2 * 3;
    int c = (2347 + C43) % 10000;
    boolean flag = false;
    string label = "rec2347";
};
type R2348 record {
    int id;
    int a = C44 + 2348;
    int b = C5 * 4;
    int c = (2348 + C44) % 10000;
    boolean flag = true;
    string label = "rec2348";
};
type R2349 record {
    int id;
    int a = C45 + 2349;
    int b = C8 * 5;
    int c = (2349 + C45) % 10000;
    boolean flag = false;
    string label = "rec2349";
};
type R2350 record {
    int id;
    int a = C46 + 2350;
    int b = C11 * 6;
    int c = (2350 + C46) % 10000;
    boolean flag = true;
    string label = "rec2350";
};
type R2351 record {
    int id;
    int a = C47 + 2351;
    int b = C14 * 7;
    int c = (2351 + C47) % 10000;
    boolean flag = false;
    string label = "rec2351";
};
type R2352 record {
    int id;
    int a = C48 + 2352;
    int b = C17 * 1;
    int c = (2352 + C48) % 10000;
    boolean flag = true;
    string label = "rec2352";
};
type R2353 record {
    int id;
    int a = C49 + 2353;
    int b = C20 * 2;
    int c = (2353 + C49) % 10000;
    boolean flag = false;
    string label = "rec2353";
};
type R2354 record {
    int id;
    int a = C50 + 2354;
    int b = C23 * 3;
    int c = (2354 + C50) % 10000;
    boolean flag = true;
    string label = "rec2354";
};
type R2355 record {
    int id;
    int a = C51 + 2355;
    int b = C26 * 4;
    int c = (2355 + C51) % 10000;
    boolean flag = false;
    string label = "rec2355";
};
type R2356 record {
    int id;
    int a = C52 + 2356;
    int b = C29 * 5;
    int c = (2356 + C52) % 10000;
    boolean flag = true;
    string label = "rec2356";
};
type R2357 record {
    int id;
    int a = C53 + 2357;
    int b = C32 * 6;
    int c = (2357 + C53) % 10000;
    boolean flag = false;
    string label = "rec2357";
};
type R2358 record {
    int id;
    int a = C54 + 2358;
    int b = C35 * 7;
    int c = (2358 + C54) % 10000;
    boolean flag = true;
    string label = "rec2358";
};
type R2359 record {
    int id;
    int a = C55 + 2359;
    int b = C38 * 1;
    int c = (2359 + C55) % 10000;
    boolean flag = false;
    string label = "rec2359";
};
type R2360 record {
    int id;
    int a = C56 + 2360;
    int b = C41 * 2;
    int c = (2360 + C56) % 10000;
    boolean flag = true;
    string label = "rec2360";
};
type R2361 record {
    int id;
    int a = C57 + 2361;
    int b = C44 * 3;
    int c = (2361 + C57) % 10000;
    boolean flag = false;
    string label = "rec2361";
};
type R2362 record {
    int id;
    int a = C58 + 2362;
    int b = C47 * 4;
    int c = (2362 + C58) % 10000;
    boolean flag = true;
    string label = "rec2362";
};
type R2363 record {
    int id;
    int a = C59 + 2363;
    int b = C50 * 5;
    int c = (2363 + C59) % 10000;
    boolean flag = false;
    string label = "rec2363";
};
type R2364 record {
    int id;
    int a = C60 + 2364;
    int b = C53 * 6;
    int c = (2364 + C60) % 10000;
    boolean flag = true;
    string label = "rec2364";
};
type R2365 record {
    int id;
    int a = C61 + 2365;
    int b = C56 * 7;
    int c = (2365 + C61) % 10000;
    boolean flag = false;
    string label = "rec2365";
};
type R2366 record {
    int id;
    int a = C62 + 2366;
    int b = C59 * 1;
    int c = (2366 + C62) % 10000;
    boolean flag = true;
    string label = "rec2366";
};
type R2367 record {
    int id;
    int a = C63 + 2367;
    int b = C62 * 2;
    int c = (2367 + C63) % 10000;
    boolean flag = false;
    string label = "rec2367";
};
type R2368 record {
    int id;
    int a = C0 + 2368;
    int b = C1 * 3;
    int c = (2368 + C0) % 10000;
    boolean flag = true;
    string label = "rec2368";
};
type R2369 record {
    int id;
    int a = C1 + 2369;
    int b = C4 * 4;
    int c = (2369 + C1) % 10000;
    boolean flag = false;
    string label = "rec2369";
};
type R2370 record {
    int id;
    int a = C2 + 2370;
    int b = C7 * 5;
    int c = (2370 + C2) % 10000;
    boolean flag = true;
    string label = "rec2370";
};
type R2371 record {
    int id;
    int a = C3 + 2371;
    int b = C10 * 6;
    int c = (2371 + C3) % 10000;
    boolean flag = false;
    string label = "rec2371";
};
type R2372 record {
    int id;
    int a = C4 + 2372;
    int b = C13 * 7;
    int c = (2372 + C4) % 10000;
    boolean flag = true;
    string label = "rec2372";
};
type R2373 record {
    int id;
    int a = C5 + 2373;
    int b = C16 * 1;
    int c = (2373 + C5) % 10000;
    boolean flag = false;
    string label = "rec2373";
};
type R2374 record {
    int id;
    int a = C6 + 2374;
    int b = C19 * 2;
    int c = (2374 + C6) % 10000;
    boolean flag = true;
    string label = "rec2374";
};
type R2375 record {
    int id;
    int a = C7 + 2375;
    int b = C22 * 3;
    int c = (2375 + C7) % 10000;
    boolean flag = false;
    string label = "rec2375";
};
type R2376 record {
    int id;
    int a = C8 + 2376;
    int b = C25 * 4;
    int c = (2376 + C8) % 10000;
    boolean flag = true;
    string label = "rec2376";
};
type R2377 record {
    int id;
    int a = C9 + 2377;
    int b = C28 * 5;
    int c = (2377 + C9) % 10000;
    boolean flag = false;
    string label = "rec2377";
};
type R2378 record {
    int id;
    int a = C10 + 2378;
    int b = C31 * 6;
    int c = (2378 + C10) % 10000;
    boolean flag = true;
    string label = "rec2378";
};
type R2379 record {
    int id;
    int a = C11 + 2379;
    int b = C34 * 7;
    int c = (2379 + C11) % 10000;
    boolean flag = false;
    string label = "rec2379";
};
type R2380 record {
    int id;
    int a = C12 + 2380;
    int b = C37 * 1;
    int c = (2380 + C12) % 10000;
    boolean flag = true;
    string label = "rec2380";
};
type R2381 record {
    int id;
    int a = C13 + 2381;
    int b = C40 * 2;
    int c = (2381 + C13) % 10000;
    boolean flag = false;
    string label = "rec2381";
};
type R2382 record {
    int id;
    int a = C14 + 2382;
    int b = C43 * 3;
    int c = (2382 + C14) % 10000;
    boolean flag = true;
    string label = "rec2382";
};
type R2383 record {
    int id;
    int a = C15 + 2383;
    int b = C46 * 4;
    int c = (2383 + C15) % 10000;
    boolean flag = false;
    string label = "rec2383";
};
type R2384 record {
    int id;
    int a = C16 + 2384;
    int b = C49 * 5;
    int c = (2384 + C16) % 10000;
    boolean flag = true;
    string label = "rec2384";
};
type R2385 record {
    int id;
    int a = C17 + 2385;
    int b = C52 * 6;
    int c = (2385 + C17) % 10000;
    boolean flag = false;
    string label = "rec2385";
};
type R2386 record {
    int id;
    int a = C18 + 2386;
    int b = C55 * 7;
    int c = (2386 + C18) % 10000;
    boolean flag = true;
    string label = "rec2386";
};
type R2387 record {
    int id;
    int a = C19 + 2387;
    int b = C58 * 1;
    int c = (2387 + C19) % 10000;
    boolean flag = false;
    string label = "rec2387";
};
type R2388 record {
    int id;
    int a = C20 + 2388;
    int b = C61 * 2;
    int c = (2388 + C20) % 10000;
    boolean flag = true;
    string label = "rec2388";
};
type R2389 record {
    int id;
    int a = C21 + 2389;
    int b = C0 * 3;
    int c = (2389 + C21) % 10000;
    boolean flag = false;
    string label = "rec2389";
};
type R2390 record {
    int id;
    int a = C22 + 2390;
    int b = C3 * 4;
    int c = (2390 + C22) % 10000;
    boolean flag = true;
    string label = "rec2390";
};
type R2391 record {
    int id;
    int a = C23 + 2391;
    int b = C6 * 5;
    int c = (2391 + C23) % 10000;
    boolean flag = false;
    string label = "rec2391";
};
type R2392 record {
    int id;
    int a = C24 + 2392;
    int b = C9 * 6;
    int c = (2392 + C24) % 10000;
    boolean flag = true;
    string label = "rec2392";
};
type R2393 record {
    int id;
    int a = C25 + 2393;
    int b = C12 * 7;
    int c = (2393 + C25) % 10000;
    boolean flag = false;
    string label = "rec2393";
};
type R2394 record {
    int id;
    int a = C26 + 2394;
    int b = C15 * 1;
    int c = (2394 + C26) % 10000;
    boolean flag = true;
    string label = "rec2394";
};
type R2395 record {
    int id;
    int a = C27 + 2395;
    int b = C18 * 2;
    int c = (2395 + C27) % 10000;
    boolean flag = false;
    string label = "rec2395";
};
type R2396 record {
    int id;
    int a = C28 + 2396;
    int b = C21 * 3;
    int c = (2396 + C28) % 10000;
    boolean flag = true;
    string label = "rec2396";
};
type R2397 record {
    int id;
    int a = C29 + 2397;
    int b = C24 * 4;
    int c = (2397 + C29) % 10000;
    boolean flag = false;
    string label = "rec2397";
};
type R2398 record {
    int id;
    int a = C30 + 2398;
    int b = C27 * 5;
    int c = (2398 + C30) % 10000;
    boolean flag = true;
    string label = "rec2398";
};
type R2399 record {
    int id;
    int a = C31 + 2399;
    int b = C30 * 6;
    int c = (2399 + C31) % 10000;
    boolean flag = false;
    string label = "rec2399";
};
type R2400 record {
    int id;
    int a = C32 + 2400;
    int b = C33 * 7;
    int c = (2400 + C32) % 10000;
    boolean flag = true;
    string label = "rec2400";
};
type R2401 record {
    int id;
    int a = C33 + 2401;
    int b = C36 * 1;
    int c = (2401 + C33) % 10000;
    boolean flag = false;
    string label = "rec2401";
};
type R2402 record {
    int id;
    int a = C34 + 2402;
    int b = C39 * 2;
    int c = (2402 + C34) % 10000;
    boolean flag = true;
    string label = "rec2402";
};
type R2403 record {
    int id;
    int a = C35 + 2403;
    int b = C42 * 3;
    int c = (2403 + C35) % 10000;
    boolean flag = false;
    string label = "rec2403";
};
type R2404 record {
    int id;
    int a = C36 + 2404;
    int b = C45 * 4;
    int c = (2404 + C36) % 10000;
    boolean flag = true;
    string label = "rec2404";
};
type R2405 record {
    int id;
    int a = C37 + 2405;
    int b = C48 * 5;
    int c = (2405 + C37) % 10000;
    boolean flag = false;
    string label = "rec2405";
};
type R2406 record {
    int id;
    int a = C38 + 2406;
    int b = C51 * 6;
    int c = (2406 + C38) % 10000;
    boolean flag = true;
    string label = "rec2406";
};
type R2407 record {
    int id;
    int a = C39 + 2407;
    int b = C54 * 7;
    int c = (2407 + C39) % 10000;
    boolean flag = false;
    string label = "rec2407";
};
type R2408 record {
    int id;
    int a = C40 + 2408;
    int b = C57 * 1;
    int c = (2408 + C40) % 10000;
    boolean flag = true;
    string label = "rec2408";
};
type R2409 record {
    int id;
    int a = C41 + 2409;
    int b = C60 * 2;
    int c = (2409 + C41) % 10000;
    boolean flag = false;
    string label = "rec2409";
};
type R2410 record {
    int id;
    int a = C42 + 2410;
    int b = C63 * 3;
    int c = (2410 + C42) % 10000;
    boolean flag = true;
    string label = "rec2410";
};
type R2411 record {
    int id;
    int a = C43 + 2411;
    int b = C2 * 4;
    int c = (2411 + C43) % 10000;
    boolean flag = false;
    string label = "rec2411";
};
type R2412 record {
    int id;
    int a = C44 + 2412;
    int b = C5 * 5;
    int c = (2412 + C44) % 10000;
    boolean flag = true;
    string label = "rec2412";
};
type R2413 record {
    int id;
    int a = C45 + 2413;
    int b = C8 * 6;
    int c = (2413 + C45) % 10000;
    boolean flag = false;
    string label = "rec2413";
};
type R2414 record {
    int id;
    int a = C46 + 2414;
    int b = C11 * 7;
    int c = (2414 + C46) % 10000;
    boolean flag = true;
    string label = "rec2414";
};
type R2415 record {
    int id;
    int a = C47 + 2415;
    int b = C14 * 1;
    int c = (2415 + C47) % 10000;
    boolean flag = false;
    string label = "rec2415";
};
type R2416 record {
    int id;
    int a = C48 + 2416;
    int b = C17 * 2;
    int c = (2416 + C48) % 10000;
    boolean flag = true;
    string label = "rec2416";
};
type R2417 record {
    int id;
    int a = C49 + 2417;
    int b = C20 * 3;
    int c = (2417 + C49) % 10000;
    boolean flag = false;
    string label = "rec2417";
};
type R2418 record {
    int id;
    int a = C50 + 2418;
    int b = C23 * 4;
    int c = (2418 + C50) % 10000;
    boolean flag = true;
    string label = "rec2418";
};
type R2419 record {
    int id;
    int a = C51 + 2419;
    int b = C26 * 5;
    int c = (2419 + C51) % 10000;
    boolean flag = false;
    string label = "rec2419";
};
type R2420 record {
    int id;
    int a = C52 + 2420;
    int b = C29 * 6;
    int c = (2420 + C52) % 10000;
    boolean flag = true;
    string label = "rec2420";
};
type R2421 record {
    int id;
    int a = C53 + 2421;
    int b = C32 * 7;
    int c = (2421 + C53) % 10000;
    boolean flag = false;
    string label = "rec2421";
};
type R2422 record {
    int id;
    int a = C54 + 2422;
    int b = C35 * 1;
    int c = (2422 + C54) % 10000;
    boolean flag = true;
    string label = "rec2422";
};
type R2423 record {
    int id;
    int a = C55 + 2423;
    int b = C38 * 2;
    int c = (2423 + C55) % 10000;
    boolean flag = false;
    string label = "rec2423";
};
type R2424 record {
    int id;
    int a = C56 + 2424;
    int b = C41 * 3;
    int c = (2424 + C56) % 10000;
    boolean flag = true;
    string label = "rec2424";
};
type R2425 record {
    int id;
    int a = C57 + 2425;
    int b = C44 * 4;
    int c = (2425 + C57) % 10000;
    boolean flag = false;
    string label = "rec2425";
};
type R2426 record {
    int id;
    int a = C58 + 2426;
    int b = C47 * 5;
    int c = (2426 + C58) % 10000;
    boolean flag = true;
    string label = "rec2426";
};
type R2427 record {
    int id;
    int a = C59 + 2427;
    int b = C50 * 6;
    int c = (2427 + C59) % 10000;
    boolean flag = false;
    string label = "rec2427";
};
type R2428 record {
    int id;
    int a = C60 + 2428;
    int b = C53 * 7;
    int c = (2428 + C60) % 10000;
    boolean flag = true;
    string label = "rec2428";
};
type R2429 record {
    int id;
    int a = C61 + 2429;
    int b = C56 * 1;
    int c = (2429 + C61) % 10000;
    boolean flag = false;
    string label = "rec2429";
};
type R2430 record {
    int id;
    int a = C62 + 2430;
    int b = C59 * 2;
    int c = (2430 + C62) % 10000;
    boolean flag = true;
    string label = "rec2430";
};
type R2431 record {
    int id;
    int a = C63 + 2431;
    int b = C62 * 3;
    int c = (2431 + C63) % 10000;
    boolean flag = false;
    string label = "rec2431";
};
type R2432 record {
    int id;
    int a = C0 + 2432;
    int b = C1 * 4;
    int c = (2432 + C0) % 10000;
    boolean flag = true;
    string label = "rec2432";
};
type R2433 record {
    int id;
    int a = C1 + 2433;
    int b = C4 * 5;
    int c = (2433 + C1) % 10000;
    boolean flag = false;
    string label = "rec2433";
};
type R2434 record {
    int id;
    int a = C2 + 2434;
    int b = C7 * 6;
    int c = (2434 + C2) % 10000;
    boolean flag = true;
    string label = "rec2434";
};
type R2435 record {
    int id;
    int a = C3 + 2435;
    int b = C10 * 7;
    int c = (2435 + C3) % 10000;
    boolean flag = false;
    string label = "rec2435";
};
type R2436 record {
    int id;
    int a = C4 + 2436;
    int b = C13 * 1;
    int c = (2436 + C4) % 10000;
    boolean flag = true;
    string label = "rec2436";
};
type R2437 record {
    int id;
    int a = C5 + 2437;
    int b = C16 * 2;
    int c = (2437 + C5) % 10000;
    boolean flag = false;
    string label = "rec2437";
};
type R2438 record {
    int id;
    int a = C6 + 2438;
    int b = C19 * 3;
    int c = (2438 + C6) % 10000;
    boolean flag = true;
    string label = "rec2438";
};
type R2439 record {
    int id;
    int a = C7 + 2439;
    int b = C22 * 4;
    int c = (2439 + C7) % 10000;
    boolean flag = false;
    string label = "rec2439";
};
type R2440 record {
    int id;
    int a = C8 + 2440;
    int b = C25 * 5;
    int c = (2440 + C8) % 10000;
    boolean flag = true;
    string label = "rec2440";
};
type R2441 record {
    int id;
    int a = C9 + 2441;
    int b = C28 * 6;
    int c = (2441 + C9) % 10000;
    boolean flag = false;
    string label = "rec2441";
};
type R2442 record {
    int id;
    int a = C10 + 2442;
    int b = C31 * 7;
    int c = (2442 + C10) % 10000;
    boolean flag = true;
    string label = "rec2442";
};
type R2443 record {
    int id;
    int a = C11 + 2443;
    int b = C34 * 1;
    int c = (2443 + C11) % 10000;
    boolean flag = false;
    string label = "rec2443";
};
type R2444 record {
    int id;
    int a = C12 + 2444;
    int b = C37 * 2;
    int c = (2444 + C12) % 10000;
    boolean flag = true;
    string label = "rec2444";
};
type R2445 record {
    int id;
    int a = C13 + 2445;
    int b = C40 * 3;
    int c = (2445 + C13) % 10000;
    boolean flag = false;
    string label = "rec2445";
};
type R2446 record {
    int id;
    int a = C14 + 2446;
    int b = C43 * 4;
    int c = (2446 + C14) % 10000;
    boolean flag = true;
    string label = "rec2446";
};
type R2447 record {
    int id;
    int a = C15 + 2447;
    int b = C46 * 5;
    int c = (2447 + C15) % 10000;
    boolean flag = false;
    string label = "rec2447";
};
type R2448 record {
    int id;
    int a = C16 + 2448;
    int b = C49 * 6;
    int c = (2448 + C16) % 10000;
    boolean flag = true;
    string label = "rec2448";
};
type R2449 record {
    int id;
    int a = C17 + 2449;
    int b = C52 * 7;
    int c = (2449 + C17) % 10000;
    boolean flag = false;
    string label = "rec2449";
};
type R2450 record {
    int id;
    int a = C18 + 2450;
    int b = C55 * 1;
    int c = (2450 + C18) % 10000;
    boolean flag = true;
    string label = "rec2450";
};
type R2451 record {
    int id;
    int a = C19 + 2451;
    int b = C58 * 2;
    int c = (2451 + C19) % 10000;
    boolean flag = false;
    string label = "rec2451";
};
type R2452 record {
    int id;
    int a = C20 + 2452;
    int b = C61 * 3;
    int c = (2452 + C20) % 10000;
    boolean flag = true;
    string label = "rec2452";
};
type R2453 record {
    int id;
    int a = C21 + 2453;
    int b = C0 * 4;
    int c = (2453 + C21) % 10000;
    boolean flag = false;
    string label = "rec2453";
};
type R2454 record {
    int id;
    int a = C22 + 2454;
    int b = C3 * 5;
    int c = (2454 + C22) % 10000;
    boolean flag = true;
    string label = "rec2454";
};
type R2455 record {
    int id;
    int a = C23 + 2455;
    int b = C6 * 6;
    int c = (2455 + C23) % 10000;
    boolean flag = false;
    string label = "rec2455";
};
type R2456 record {
    int id;
    int a = C24 + 2456;
    int b = C9 * 7;
    int c = (2456 + C24) % 10000;
    boolean flag = true;
    string label = "rec2456";
};
type R2457 record {
    int id;
    int a = C25 + 2457;
    int b = C12 * 1;
    int c = (2457 + C25) % 10000;
    boolean flag = false;
    string label = "rec2457";
};
type R2458 record {
    int id;
    int a = C26 + 2458;
    int b = C15 * 2;
    int c = (2458 + C26) % 10000;
    boolean flag = true;
    string label = "rec2458";
};
type R2459 record {
    int id;
    int a = C27 + 2459;
    int b = C18 * 3;
    int c = (2459 + C27) % 10000;
    boolean flag = false;
    string label = "rec2459";
};
type R2460 record {
    int id;
    int a = C28 + 2460;
    int b = C21 * 4;
    int c = (2460 + C28) % 10000;
    boolean flag = true;
    string label = "rec2460";
};
type R2461 record {
    int id;
    int a = C29 + 2461;
    int b = C24 * 5;
    int c = (2461 + C29) % 10000;
    boolean flag = false;
    string label = "rec2461";
};
type R2462 record {
    int id;
    int a = C30 + 2462;
    int b = C27 * 6;
    int c = (2462 + C30) % 10000;
    boolean flag = true;
    string label = "rec2462";
};
type R2463 record {
    int id;
    int a = C31 + 2463;
    int b = C30 * 7;
    int c = (2463 + C31) % 10000;
    boolean flag = false;
    string label = "rec2463";
};
type R2464 record {
    int id;
    int a = C32 + 2464;
    int b = C33 * 1;
    int c = (2464 + C32) % 10000;
    boolean flag = true;
    string label = "rec2464";
};
type R2465 record {
    int id;
    int a = C33 + 2465;
    int b = C36 * 2;
    int c = (2465 + C33) % 10000;
    boolean flag = false;
    string label = "rec2465";
};
type R2466 record {
    int id;
    int a = C34 + 2466;
    int b = C39 * 3;
    int c = (2466 + C34) % 10000;
    boolean flag = true;
    string label = "rec2466";
};
type R2467 record {
    int id;
    int a = C35 + 2467;
    int b = C42 * 4;
    int c = (2467 + C35) % 10000;
    boolean flag = false;
    string label = "rec2467";
};
type R2468 record {
    int id;
    int a = C36 + 2468;
    int b = C45 * 5;
    int c = (2468 + C36) % 10000;
    boolean flag = true;
    string label = "rec2468";
};
type R2469 record {
    int id;
    int a = C37 + 2469;
    int b = C48 * 6;
    int c = (2469 + C37) % 10000;
    boolean flag = false;
    string label = "rec2469";
};
type R2470 record {
    int id;
    int a = C38 + 2470;
    int b = C51 * 7;
    int c = (2470 + C38) % 10000;
    boolean flag = true;
    string label = "rec2470";
};
type R2471 record {
    int id;
    int a = C39 + 2471;
    int b = C54 * 1;
    int c = (2471 + C39) % 10000;
    boolean flag = false;
    string label = "rec2471";
};
type R2472 record {
    int id;
    int a = C40 + 2472;
    int b = C57 * 2;
    int c = (2472 + C40) % 10000;
    boolean flag = true;
    string label = "rec2472";
};
type R2473 record {
    int id;
    int a = C41 + 2473;
    int b = C60 * 3;
    int c = (2473 + C41) % 10000;
    boolean flag = false;
    string label = "rec2473";
};
type R2474 record {
    int id;
    int a = C42 + 2474;
    int b = C63 * 4;
    int c = (2474 + C42) % 10000;
    boolean flag = true;
    string label = "rec2474";
};
type R2475 record {
    int id;
    int a = C43 + 2475;
    int b = C2 * 5;
    int c = (2475 + C43) % 10000;
    boolean flag = false;
    string label = "rec2475";
};
type R2476 record {
    int id;
    int a = C44 + 2476;
    int b = C5 * 6;
    int c = (2476 + C44) % 10000;
    boolean flag = true;
    string label = "rec2476";
};
type R2477 record {
    int id;
    int a = C45 + 2477;
    int b = C8 * 7;
    int c = (2477 + C45) % 10000;
    boolean flag = false;
    string label = "rec2477";
};
type R2478 record {
    int id;
    int a = C46 + 2478;
    int b = C11 * 1;
    int c = (2478 + C46) % 10000;
    boolean flag = true;
    string label = "rec2478";
};
type R2479 record {
    int id;
    int a = C47 + 2479;
    int b = C14 * 2;
    int c = (2479 + C47) % 10000;
    boolean flag = false;
    string label = "rec2479";
};
type R2480 record {
    int id;
    int a = C48 + 2480;
    int b = C17 * 3;
    int c = (2480 + C48) % 10000;
    boolean flag = true;
    string label = "rec2480";
};
type R2481 record {
    int id;
    int a = C49 + 2481;
    int b = C20 * 4;
    int c = (2481 + C49) % 10000;
    boolean flag = false;
    string label = "rec2481";
};
type R2482 record {
    int id;
    int a = C50 + 2482;
    int b = C23 * 5;
    int c = (2482 + C50) % 10000;
    boolean flag = true;
    string label = "rec2482";
};
type R2483 record {
    int id;
    int a = C51 + 2483;
    int b = C26 * 6;
    int c = (2483 + C51) % 10000;
    boolean flag = false;
    string label = "rec2483";
};
type R2484 record {
    int id;
    int a = C52 + 2484;
    int b = C29 * 7;
    int c = (2484 + C52) % 10000;
    boolean flag = true;
    string label = "rec2484";
};
type R2485 record {
    int id;
    int a = C53 + 2485;
    int b = C32 * 1;
    int c = (2485 + C53) % 10000;
    boolean flag = false;
    string label = "rec2485";
};
type R2486 record {
    int id;
    int a = C54 + 2486;
    int b = C35 * 2;
    int c = (2486 + C54) % 10000;
    boolean flag = true;
    string label = "rec2486";
};
type R2487 record {
    int id;
    int a = C55 + 2487;
    int b = C38 * 3;
    int c = (2487 + C55) % 10000;
    boolean flag = false;
    string label = "rec2487";
};
type R2488 record {
    int id;
    int a = C56 + 2488;
    int b = C41 * 4;
    int c = (2488 + C56) % 10000;
    boolean flag = true;
    string label = "rec2488";
};
type R2489 record {
    int id;
    int a = C57 + 2489;
    int b = C44 * 5;
    int c = (2489 + C57) % 10000;
    boolean flag = false;
    string label = "rec2489";
};
type R2490 record {
    int id;
    int a = C58 + 2490;
    int b = C47 * 6;
    int c = (2490 + C58) % 10000;
    boolean flag = true;
    string label = "rec2490";
};
type R2491 record {
    int id;
    int a = C59 + 2491;
    int b = C50 * 7;
    int c = (2491 + C59) % 10000;
    boolean flag = false;
    string label = "rec2491";
};
type R2492 record {
    int id;
    int a = C60 + 2492;
    int b = C53 * 1;
    int c = (2492 + C60) % 10000;
    boolean flag = true;
    string label = "rec2492";
};
type R2493 record {
    int id;
    int a = C61 + 2493;
    int b = C56 * 2;
    int c = (2493 + C61) % 10000;
    boolean flag = false;
    string label = "rec2493";
};
type R2494 record {
    int id;
    int a = C62 + 2494;
    int b = C59 * 3;
    int c = (2494 + C62) % 10000;
    boolean flag = true;
    string label = "rec2494";
};
type R2495 record {
    int id;
    int a = C63 + 2495;
    int b = C62 * 4;
    int c = (2495 + C63) % 10000;
    boolean flag = false;
    string label = "rec2495";
};
type R2496 record {
    int id;
    int a = C0 + 2496;
    int b = C1 * 5;
    int c = (2496 + C0) % 10000;
    boolean flag = true;
    string label = "rec2496";
};
type R2497 record {
    int id;
    int a = C1 + 2497;
    int b = C4 * 6;
    int c = (2497 + C1) % 10000;
    boolean flag = false;
    string label = "rec2497";
};
type R2498 record {
    int id;
    int a = C2 + 2498;
    int b = C7 * 7;
    int c = (2498 + C2) % 10000;
    boolean flag = true;
    string label = "rec2498";
};
type R2499 record {
    int id;
    int a = C3 + 2499;
    int b = C10 * 1;
    int c = (2499 + C3) % 10000;
    boolean flag = false;
    string label = "rec2499";
};

function makeR0() returns R0 {
    R0 r = {id: 0};
    return r;
}
function makeR1() returns R1 {
    R1 r = {id: 1};
    return r;
}
function makeR2() returns R2 {
    R2 r = {id: 2};
    return r;
}
function makeR3() returns R3 {
    R3 r = {id: 3};
    return r;
}
function makeR4() returns R4 {
    R4 r = {id: 4};
    return r;
}
function makeR5() returns R5 {
    R5 r = {id: 5};
    return r;
}
function makeR6() returns R6 {
    R6 r = {id: 6};
    return r;
}
function makeR7() returns R7 {
    R7 r = {id: 7};
    return r;
}
function makeR8() returns R8 {
    R8 r = {id: 8};
    return r;
}
function makeR9() returns R9 {
    R9 r = {id: 9};
    return r;
}
function makeR10() returns R10 {
    R10 r = {id: 10};
    return r;
}
function makeR11() returns R11 {
    R11 r = {id: 11};
    return r;
}
function makeR12() returns R12 {
    R12 r = {id: 12};
    return r;
}
function makeR13() returns R13 {
    R13 r = {id: 13};
    return r;
}
function makeR14() returns R14 {
    R14 r = {id: 14};
    return r;
}
function makeR15() returns R15 {
    R15 r = {id: 15};
    return r;
}
function makeR16() returns R16 {
    R16 r = {id: 16};
    return r;
}
function makeR17() returns R17 {
    R17 r = {id: 17};
    return r;
}
function makeR18() returns R18 {
    R18 r = {id: 18};
    return r;
}
function makeR19() returns R19 {
    R19 r = {id: 19};
    return r;
}
function makeR20() returns R20 {
    R20 r = {id: 20};
    return r;
}
function makeR21() returns R21 {
    R21 r = {id: 21};
    return r;
}
function makeR22() returns R22 {
    R22 r = {id: 22};
    return r;
}
function makeR23() returns R23 {
    R23 r = {id: 23};
    return r;
}
function makeR24() returns R24 {
    R24 r = {id: 24};
    return r;
}
function makeR25() returns R25 {
    R25 r = {id: 25};
    return r;
}
function makeR26() returns R26 {
    R26 r = {id: 26};
    return r;
}
function makeR27() returns R27 {
    R27 r = {id: 27};
    return r;
}
function makeR28() returns R28 {
    R28 r = {id: 28};
    return r;
}
function makeR29() returns R29 {
    R29 r = {id: 29};
    return r;
}
function makeR30() returns R30 {
    R30 r = {id: 30};
    return r;
}
function makeR31() returns R31 {
    R31 r = {id: 31};
    return r;
}
function makeR32() returns R32 {
    R32 r = {id: 32};
    return r;
}
function makeR33() returns R33 {
    R33 r = {id: 33};
    return r;
}
function makeR34() returns R34 {
    R34 r = {id: 34};
    return r;
}
function makeR35() returns R35 {
    R35 r = {id: 35};
    return r;
}
function makeR36() returns R36 {
    R36 r = {id: 36};
    return r;
}
function makeR37() returns R37 {
    R37 r = {id: 37};
    return r;
}
function makeR38() returns R38 {
    R38 r = {id: 38};
    return r;
}
function makeR39() returns R39 {
    R39 r = {id: 39};
    return r;
}
function makeR40() returns R40 {
    R40 r = {id: 40};
    return r;
}
function makeR41() returns R41 {
    R41 r = {id: 41};
    return r;
}
function makeR42() returns R42 {
    R42 r = {id: 42};
    return r;
}
function makeR43() returns R43 {
    R43 r = {id: 43};
    return r;
}
function makeR44() returns R44 {
    R44 r = {id: 44};
    return r;
}
function makeR45() returns R45 {
    R45 r = {id: 45};
    return r;
}
function makeR46() returns R46 {
    R46 r = {id: 46};
    return r;
}
function makeR47() returns R47 {
    R47 r = {id: 47};
    return r;
}
function makeR48() returns R48 {
    R48 r = {id: 48};
    return r;
}
function makeR49() returns R49 {
    R49 r = {id: 49};
    return r;
}
function makeR50() returns R50 {
    R50 r = {id: 50};
    return r;
}
function makeR51() returns R51 {
    R51 r = {id: 51};
    return r;
}
function makeR52() returns R52 {
    R52 r = {id: 52};
    return r;
}
function makeR53() returns R53 {
    R53 r = {id: 53};
    return r;
}
function makeR54() returns R54 {
    R54 r = {id: 54};
    return r;
}
function makeR55() returns R55 {
    R55 r = {id: 55};
    return r;
}
function makeR56() returns R56 {
    R56 r = {id: 56};
    return r;
}
function makeR57() returns R57 {
    R57 r = {id: 57};
    return r;
}
function makeR58() returns R58 {
    R58 r = {id: 58};
    return r;
}
function makeR59() returns R59 {
    R59 r = {id: 59};
    return r;
}
function makeR60() returns R60 {
    R60 r = {id: 60};
    return r;
}
function makeR61() returns R61 {
    R61 r = {id: 61};
    return r;
}
function makeR62() returns R62 {
    R62 r = {id: 62};
    return r;
}
function makeR63() returns R63 {
    R63 r = {id: 63};
    return r;
}
function makeR64() returns R64 {
    R64 r = {id: 64};
    return r;
}
function makeR65() returns R65 {
    R65 r = {id: 65};
    return r;
}
function makeR66() returns R66 {
    R66 r = {id: 66};
    return r;
}
function makeR67() returns R67 {
    R67 r = {id: 67};
    return r;
}
function makeR68() returns R68 {
    R68 r = {id: 68};
    return r;
}
function makeR69() returns R69 {
    R69 r = {id: 69};
    return r;
}
function makeR70() returns R70 {
    R70 r = {id: 70};
    return r;
}
function makeR71() returns R71 {
    R71 r = {id: 71};
    return r;
}
function makeR72() returns R72 {
    R72 r = {id: 72};
    return r;
}
function makeR73() returns R73 {
    R73 r = {id: 73};
    return r;
}
function makeR74() returns R74 {
    R74 r = {id: 74};
    return r;
}
function makeR75() returns R75 {
    R75 r = {id: 75};
    return r;
}
function makeR76() returns R76 {
    R76 r = {id: 76};
    return r;
}
function makeR77() returns R77 {
    R77 r = {id: 77};
    return r;
}
function makeR78() returns R78 {
    R78 r = {id: 78};
    return r;
}
function makeR79() returns R79 {
    R79 r = {id: 79};
    return r;
}

public function main() {
    int total = 0;
    R0 r0 = makeR0();
    total = (total + r0.a + r0.b + r0.c) % 1000000;
    R1 r1 = makeR1();
    total = (total + r1.a + r1.b + r1.c) % 1000000;
    R2 r2 = makeR2();
    total = (total + r2.a + r2.b + r2.c) % 1000000;
    R3 r3 = makeR3();
    total = (total + r3.a + r3.b + r3.c) % 1000000;
    R4 r4 = makeR4();
    total = (total + r4.a + r4.b + r4.c) % 1000000;
    R5 r5 = makeR5();
    total = (total + r5.a + r5.b + r5.c) % 1000000;
    R6 r6 = makeR6();
    total = (total + r6.a + r6.b + r6.c) % 1000000;
    R7 r7 = makeR7();
    total = (total + r7.a + r7.b + r7.c) % 1000000;
    R8 r8 = makeR8();
    total = (total + r8.a + r8.b + r8.c) % 1000000;
    R9 r9 = makeR9();
    total = (total + r9.a + r9.b + r9.c) % 1000000;
    R10 r10 = makeR10();
    total = (total + r10.a + r10.b + r10.c) % 1000000;
    R11 r11 = makeR11();
    total = (total + r11.a + r11.b + r11.c) % 1000000;
    R12 r12 = makeR12();
    total = (total + r12.a + r12.b + r12.c) % 1000000;
    R13 r13 = makeR13();
    total = (total + r13.a + r13.b + r13.c) % 1000000;
    R14 r14 = makeR14();
    total = (total + r14.a + r14.b + r14.c) % 1000000;
    R15 r15 = makeR15();
    total = (total + r15.a + r15.b + r15.c) % 1000000;
    R16 r16 = makeR16();
    total = (total + r16.a + r16.b + r16.c) % 1000000;
    R17 r17 = makeR17();
    total = (total + r17.a + r17.b + r17.c) % 1000000;
    R18 r18 = makeR18();
    total = (total + r18.a + r18.b + r18.c) % 1000000;
    R19 r19 = makeR19();
    total = (total + r19.a + r19.b + r19.c) % 1000000;
    R20 r20 = makeR20();
    total = (total + r20.a + r20.b + r20.c) % 1000000;
    R21 r21 = makeR21();
    total = (total + r21.a + r21.b + r21.c) % 1000000;
    R22 r22 = makeR22();
    total = (total + r22.a + r22.b + r22.c) % 1000000;
    R23 r23 = makeR23();
    total = (total + r23.a + r23.b + r23.c) % 1000000;
    R24 r24 = makeR24();
    total = (total + r24.a + r24.b + r24.c) % 1000000;
    R25 r25 = makeR25();
    total = (total + r25.a + r25.b + r25.c) % 1000000;
    R26 r26 = makeR26();
    total = (total + r26.a + r26.b + r26.c) % 1000000;
    R27 r27 = makeR27();
    total = (total + r27.a + r27.b + r27.c) % 1000000;
    R28 r28 = makeR28();
    total = (total + r28.a + r28.b + r28.c) % 1000000;
    R29 r29 = makeR29();
    total = (total + r29.a + r29.b + r29.c) % 1000000;
    R30 r30 = makeR30();
    total = (total + r30.a + r30.b + r30.c) % 1000000;
    R31 r31 = makeR31();
    total = (total + r31.a + r31.b + r31.c) % 1000000;
    R32 r32 = makeR32();
    total = (total + r32.a + r32.b + r32.c) % 1000000;
    R33 r33 = makeR33();
    total = (total + r33.a + r33.b + r33.c) % 1000000;
    R34 r34 = makeR34();
    total = (total + r34.a + r34.b + r34.c) % 1000000;
    R35 r35 = makeR35();
    total = (total + r35.a + r35.b + r35.c) % 1000000;
    R36 r36 = makeR36();
    total = (total + r36.a + r36.b + r36.c) % 1000000;
    R37 r37 = makeR37();
    total = (total + r37.a + r37.b + r37.c) % 1000000;
    R38 r38 = makeR38();
    total = (total + r38.a + r38.b + r38.c) % 1000000;
    R39 r39 = makeR39();
    total = (total + r39.a + r39.b + r39.c) % 1000000;
    R40 r40 = makeR40();
    total = (total + r40.a + r40.b + r40.c) % 1000000;
    R41 r41 = makeR41();
    total = (total + r41.a + r41.b + r41.c) % 1000000;
    R42 r42 = makeR42();
    total = (total + r42.a + r42.b + r42.c) % 1000000;
    R43 r43 = makeR43();
    total = (total + r43.a + r43.b + r43.c) % 1000000;
    R44 r44 = makeR44();
    total = (total + r44.a + r44.b + r44.c) % 1000000;
    R45 r45 = makeR45();
    total = (total + r45.a + r45.b + r45.c) % 1000000;
    R46 r46 = makeR46();
    total = (total + r46.a + r46.b + r46.c) % 1000000;
    R47 r47 = makeR47();
    total = (total + r47.a + r47.b + r47.c) % 1000000;
    R48 r48 = makeR48();
    total = (total + r48.a + r48.b + r48.c) % 1000000;
    R49 r49 = makeR49();
    total = (total + r49.a + r49.b + r49.c) % 1000000;
    R50 r50 = makeR50();
    total = (total + r50.a + r50.b + r50.c) % 1000000;
    R51 r51 = makeR51();
    total = (total + r51.a + r51.b + r51.c) % 1000000;
    R52 r52 = makeR52();
    total = (total + r52.a + r52.b + r52.c) % 1000000;
    R53 r53 = makeR53();
    total = (total + r53.a + r53.b + r53.c) % 1000000;
    R54 r54 = makeR54();
    total = (total + r54.a + r54.b + r54.c) % 1000000;
    R55 r55 = makeR55();
    total = (total + r55.a + r55.b + r55.c) % 1000000;
    R56 r56 = makeR56();
    total = (total + r56.a + r56.b + r56.c) % 1000000;
    R57 r57 = makeR57();
    total = (total + r57.a + r57.b + r57.c) % 1000000;
    R58 r58 = makeR58();
    total = (total + r58.a + r58.b + r58.c) % 1000000;
    R59 r59 = makeR59();
    total = (total + r59.a + r59.b + r59.c) % 1000000;
    R60 r60 = makeR60();
    total = (total + r60.a + r60.b + r60.c) % 1000000;
    R61 r61 = makeR61();
    total = (total + r61.a + r61.b + r61.c) % 1000000;
    R62 r62 = makeR62();
    total = (total + r62.a + r62.b + r62.c) % 1000000;
    R63 r63 = makeR63();
    total = (total + r63.a + r63.b + r63.c) % 1000000;
    R64 r64 = makeR64();
    total = (total + r64.a + r64.b + r64.c) % 1000000;
    R65 r65 = makeR65();
    total = (total + r65.a + r65.b + r65.c) % 1000000;
    R66 r66 = makeR66();
    total = (total + r66.a + r66.b + r66.c) % 1000000;
    R67 r67 = makeR67();
    total = (total + r67.a + r67.b + r67.c) % 1000000;
    R68 r68 = makeR68();
    total = (total + r68.a + r68.b + r68.c) % 1000000;
    R69 r69 = makeR69();
    total = (total + r69.a + r69.b + r69.c) % 1000000;
    R70 r70 = makeR70();
    total = (total + r70.a + r70.b + r70.c) % 1000000;
    R71 r71 = makeR71();
    total = (total + r71.a + r71.b + r71.c) % 1000000;
    R72 r72 = makeR72();
    total = (total + r72.a + r72.b + r72.c) % 1000000;
    R73 r73 = makeR73();
    total = (total + r73.a + r73.b + r73.c) % 1000000;
    R74 r74 = makeR74();
    total = (total + r74.a + r74.b + r74.c) % 1000000;
    R75 r75 = makeR75();
    total = (total + r75.a + r75.b + r75.c) % 1000000;
    R76 r76 = makeR76();
    total = (total + r76.a + r76.b + r76.c) % 1000000;
    R77 r77 = makeR77();
    total = (total + r77.a + r77.b + r77.c) % 1000000;
    R78 r78 = makeR78();
    total = (total + r78.a + r78.b + r78.c) % 1000000;
    R79 r79 = makeR79();
    total = (total + r79.a + r79.b + r79.c) % 1000000;
    io:println(total);
}
