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

const int K0 = 3;
const int K1 = 10;
const int K2 = 17;
const int K3 = 24;
const int K4 = 31;
const int K5 = 38;
const int K6 = 45;
const int K7 = 52;
const int K8 = 59;
const int K9 = 66;
const int K10 = 73;
const int K11 = 80;
const int K12 = 87;
const int K13 = 94;
const int K14 = 101;
const int K15 = 108;
const int K16 = 115;
const int K17 = 122;
const int K18 = 129;
const int K19 = 136;
const int K20 = 143;
const int K21 = 150;
const int K22 = 157;
const int K23 = 164;
const int K24 = 171;
const int K25 = 178;
const int K26 = 185;
const int K27 = 192;
const int K28 = 199;
const int K29 = 206;
const int K30 = 213;
const int K31 = 220;
const int K32 = 227;
const int K33 = 234;
const int K34 = 241;
const int K35 = 248;
const int K36 = 255;
const int K37 = 262;
const int K38 = 269;
const int K39 = 276;
const int K40 = 283;
const int K41 = 290;
const int K42 = 297;
const int K43 = 304;
const int K44 = 311;
const int K45 = 318;
const int K46 = 325;
const int K47 = 332;
const int K48 = 339;
const int K49 = 346;
const int K50 = 353;
const int K51 = 360;
const int K52 = 367;
const int K53 = 374;
const int K54 = 381;
const int K55 = 388;
const int K56 = 395;
const int K57 = 402;
const int K58 = 409;
const int K59 = 416;
const int K60 = 423;
const int K61 = 430;
const int K62 = 437;
const int K63 = 444;
const int K64 = 451;
const int K65 = 458;
const int K66 = 465;
const int K67 = 472;
const int K68 = 479;
const int K69 = 486;
const int K70 = 493;
const int K71 = 500;
const int K72 = 507;
const int K73 = 514;
const int K74 = 521;
const int K75 = 528;
const int K76 = 535;
const int K77 = 542;
const int K78 = 549;
const int K79 = 556;
const int K80 = 563;
const int K81 = 570;
const int K82 = 577;
const int K83 = 584;
const int K84 = 591;
const int K85 = 598;
const int K86 = 605;
const int K87 = 612;
const int K88 = 619;
const int K89 = 626;
const int K90 = 633;
const int K91 = 640;
const int K92 = 647;
const int K93 = 654;
const int K94 = 661;
const int K95 = 668;
const int K96 = 675;
const int K97 = 682;
const int K98 = 689;
const int K99 = 696;
const int K100 = 703;
const int K101 = 710;
const int K102 = 717;
const int K103 = 724;
const int K104 = 731;
const int K105 = 738;
const int K106 = 745;
const int K107 = 752;
const int K108 = 759;
const int K109 = 766;
const int K110 = 773;
const int K111 = 780;
const int K112 = 787;
const int K113 = 794;
const int K114 = 801;
const int K115 = 808;
const int K116 = 815;
const int K117 = 822;
const int K118 = 829;
const int K119 = 836;
const int K120 = 843;
const int K121 = 850;
const int K122 = 857;
const int K123 = 864;
const int K124 = 871;
const int K125 = 878;
const int K126 = 885;
const int K127 = 892;
const int K128 = 899;
const int K129 = 906;
const int K130 = 913;
const int K131 = 920;
const int K132 = 927;
const int K133 = 934;
const int K134 = 941;
const int K135 = 948;
const int K136 = 955;
const int K137 = 962;
const int K138 = 969;
const int K139 = 976;
const int K140 = 983;
const int K141 = 990;
const int K142 = 997;
const int K143 = 1004;
const int K144 = 1011;
const int K145 = 1018;
const int K146 = 1025;
const int K147 = 1032;
const int K148 = 1039;
const int K149 = 1046;
const int K150 = 1053;
const int K151 = 1060;
const int K152 = 1067;
const int K153 = 1074;
const int K154 = 1081;
const int K155 = 1088;
const int K156 = 1095;
const int K157 = 1102;
const int K158 = 1109;
const int K159 = 1116;
const int K160 = 1123;
const int K161 = 1130;
const int K162 = 1137;
const int K163 = 1144;
const int K164 = 1151;
const int K165 = 1158;
const int K166 = 1165;
const int K167 = 1172;
const int K168 = 1179;
const int K169 = 1186;
const int K170 = 1193;
const int K171 = 1200;
const int K172 = 1207;
const int K173 = 1214;
const int K174 = 1221;
const int K175 = 1228;
const int K176 = 1235;
const int K177 = 1242;
const int K178 = 1249;
const int K179 = 1256;
const int K180 = 1263;
const int K181 = 1270;
const int K182 = 1277;
const int K183 = 1284;
const int K184 = 1291;
const int K185 = 1298;
const int K186 = 1305;
const int K187 = 1312;
const int K188 = 1319;
const int K189 = 1326;
const int K190 = 1333;
const int K191 = 1340;
const int K192 = 1347;
const int K193 = 1354;
const int K194 = 1361;
const int K195 = 1368;
const int K196 = 1375;
const int K197 = 1382;
const int K198 = 1389;
const int K199 = 1396;

class Box0 {
    int value;
    int weight;
    Box0? next;

    function init(int v) {
        self.value = v + K0;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box0 n) {
        self.next = n;
    }
}
class Box1 {
    int value;
    int weight;
    Box1? next;

    function init(int v) {
        self.value = v + K1;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box1 n) {
        self.next = n;
    }
}
class Box2 {
    int value;
    int weight;
    Box2? next;

    function init(int v) {
        self.value = v + K2;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box2 n) {
        self.next = n;
    }
}
class Box3 {
    int value;
    int weight;
    Box3? next;

    function init(int v) {
        self.value = v + K3;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box3 n) {
        self.next = n;
    }
}
class Box4 {
    int value;
    int weight;
    Box4? next;

    function init(int v) {
        self.value = v + K4;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box4 n) {
        self.next = n;
    }
}
class Box5 {
    int value;
    int weight;
    Box5? next;

    function init(int v) {
        self.value = v + K5;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box5 n) {
        self.next = n;
    }
}
class Box6 {
    int value;
    int weight;
    Box6? next;

    function init(int v) {
        self.value = v + K6;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box6 n) {
        self.next = n;
    }
}
class Box7 {
    int value;
    int weight;
    Box7? next;

    function init(int v) {
        self.value = v + K7;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box7 n) {
        self.next = n;
    }
}
class Box8 {
    int value;
    int weight;
    Box8? next;

    function init(int v) {
        self.value = v + K8;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box8 n) {
        self.next = n;
    }
}
class Box9 {
    int value;
    int weight;
    Box9? next;

    function init(int v) {
        self.value = v + K9;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box9 n) {
        self.next = n;
    }
}
class Box10 {
    int value;
    int weight;
    Box10? next;

    function init(int v) {
        self.value = v + K10;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box10 n) {
        self.next = n;
    }
}
class Box11 {
    int value;
    int weight;
    Box11? next;

    function init(int v) {
        self.value = v + K11;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box11 n) {
        self.next = n;
    }
}
class Box12 {
    int value;
    int weight;
    Box12? next;

    function init(int v) {
        self.value = v + K12;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box12 n) {
        self.next = n;
    }
}
class Box13 {
    int value;
    int weight;
    Box13? next;

    function init(int v) {
        self.value = v + K13;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box13 n) {
        self.next = n;
    }
}
class Box14 {
    int value;
    int weight;
    Box14? next;

    function init(int v) {
        self.value = v + K14;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box14 n) {
        self.next = n;
    }
}
class Box15 {
    int value;
    int weight;
    Box15? next;

    function init(int v) {
        self.value = v + K15;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box15 n) {
        self.next = n;
    }
}
class Box16 {
    int value;
    int weight;
    Box16? next;

    function init(int v) {
        self.value = v + K16;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box16 n) {
        self.next = n;
    }
}
class Box17 {
    int value;
    int weight;
    Box17? next;

    function init(int v) {
        self.value = v + K17;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box17 n) {
        self.next = n;
    }
}
class Box18 {
    int value;
    int weight;
    Box18? next;

    function init(int v) {
        self.value = v + K18;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box18 n) {
        self.next = n;
    }
}
class Box19 {
    int value;
    int weight;
    Box19? next;

    function init(int v) {
        self.value = v + K19;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box19 n) {
        self.next = n;
    }
}
class Box20 {
    int value;
    int weight;
    Box20? next;

    function init(int v) {
        self.value = v + K20;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box20 n) {
        self.next = n;
    }
}
class Box21 {
    int value;
    int weight;
    Box21? next;

    function init(int v) {
        self.value = v + K21;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box21 n) {
        self.next = n;
    }
}
class Box22 {
    int value;
    int weight;
    Box22? next;

    function init(int v) {
        self.value = v + K22;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box22 n) {
        self.next = n;
    }
}
class Box23 {
    int value;
    int weight;
    Box23? next;

    function init(int v) {
        self.value = v + K23;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box23 n) {
        self.next = n;
    }
}
class Box24 {
    int value;
    int weight;
    Box24? next;

    function init(int v) {
        self.value = v + K24;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box24 n) {
        self.next = n;
    }
}
class Box25 {
    int value;
    int weight;
    Box25? next;

    function init(int v) {
        self.value = v + K25;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box25 n) {
        self.next = n;
    }
}
class Box26 {
    int value;
    int weight;
    Box26? next;

    function init(int v) {
        self.value = v + K26;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box26 n) {
        self.next = n;
    }
}
class Box27 {
    int value;
    int weight;
    Box27? next;

    function init(int v) {
        self.value = v + K27;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box27 n) {
        self.next = n;
    }
}
class Box28 {
    int value;
    int weight;
    Box28? next;

    function init(int v) {
        self.value = v + K28;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box28 n) {
        self.next = n;
    }
}
class Box29 {
    int value;
    int weight;
    Box29? next;

    function init(int v) {
        self.value = v + K29;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box29 n) {
        self.next = n;
    }
}
class Box30 {
    int value;
    int weight;
    Box30? next;

    function init(int v) {
        self.value = v + K30;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box30 n) {
        self.next = n;
    }
}
class Box31 {
    int value;
    int weight;
    Box31? next;

    function init(int v) {
        self.value = v + K31;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box31 n) {
        self.next = n;
    }
}
class Box32 {
    int value;
    int weight;
    Box32? next;

    function init(int v) {
        self.value = v + K32;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box32 n) {
        self.next = n;
    }
}
class Box33 {
    int value;
    int weight;
    Box33? next;

    function init(int v) {
        self.value = v + K33;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box33 n) {
        self.next = n;
    }
}
class Box34 {
    int value;
    int weight;
    Box34? next;

    function init(int v) {
        self.value = v + K34;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box34 n) {
        self.next = n;
    }
}
class Box35 {
    int value;
    int weight;
    Box35? next;

    function init(int v) {
        self.value = v + K35;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box35 n) {
        self.next = n;
    }
}
class Box36 {
    int value;
    int weight;
    Box36? next;

    function init(int v) {
        self.value = v + K36;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box36 n) {
        self.next = n;
    }
}
class Box37 {
    int value;
    int weight;
    Box37? next;

    function init(int v) {
        self.value = v + K37;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box37 n) {
        self.next = n;
    }
}
class Box38 {
    int value;
    int weight;
    Box38? next;

    function init(int v) {
        self.value = v + K38;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box38 n) {
        self.next = n;
    }
}
class Box39 {
    int value;
    int weight;
    Box39? next;

    function init(int v) {
        self.value = v + K39;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box39 n) {
        self.next = n;
    }
}
class Box40 {
    int value;
    int weight;
    Box40? next;

    function init(int v) {
        self.value = v + K40;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box40 n) {
        self.next = n;
    }
}
class Box41 {
    int value;
    int weight;
    Box41? next;

    function init(int v) {
        self.value = v + K41;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box41 n) {
        self.next = n;
    }
}
class Box42 {
    int value;
    int weight;
    Box42? next;

    function init(int v) {
        self.value = v + K42;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box42 n) {
        self.next = n;
    }
}
class Box43 {
    int value;
    int weight;
    Box43? next;

    function init(int v) {
        self.value = v + K43;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box43 n) {
        self.next = n;
    }
}
class Box44 {
    int value;
    int weight;
    Box44? next;

    function init(int v) {
        self.value = v + K44;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box44 n) {
        self.next = n;
    }
}
class Box45 {
    int value;
    int weight;
    Box45? next;

    function init(int v) {
        self.value = v + K45;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box45 n) {
        self.next = n;
    }
}
class Box46 {
    int value;
    int weight;
    Box46? next;

    function init(int v) {
        self.value = v + K46;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box46 n) {
        self.next = n;
    }
}
class Box47 {
    int value;
    int weight;
    Box47? next;

    function init(int v) {
        self.value = v + K47;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box47 n) {
        self.next = n;
    }
}
class Box48 {
    int value;
    int weight;
    Box48? next;

    function init(int v) {
        self.value = v + K48;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box48 n) {
        self.next = n;
    }
}
class Box49 {
    int value;
    int weight;
    Box49? next;

    function init(int v) {
        self.value = v + K49;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box49 n) {
        self.next = n;
    }
}
class Box50 {
    int value;
    int weight;
    Box50? next;

    function init(int v) {
        self.value = v + K50;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box50 n) {
        self.next = n;
    }
}
class Box51 {
    int value;
    int weight;
    Box51? next;

    function init(int v) {
        self.value = v + K51;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box51 n) {
        self.next = n;
    }
}
class Box52 {
    int value;
    int weight;
    Box52? next;

    function init(int v) {
        self.value = v + K52;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box52 n) {
        self.next = n;
    }
}
class Box53 {
    int value;
    int weight;
    Box53? next;

    function init(int v) {
        self.value = v + K53;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box53 n) {
        self.next = n;
    }
}
class Box54 {
    int value;
    int weight;
    Box54? next;

    function init(int v) {
        self.value = v + K54;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box54 n) {
        self.next = n;
    }
}
class Box55 {
    int value;
    int weight;
    Box55? next;

    function init(int v) {
        self.value = v + K55;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box55 n) {
        self.next = n;
    }
}
class Box56 {
    int value;
    int weight;
    Box56? next;

    function init(int v) {
        self.value = v + K56;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box56 n) {
        self.next = n;
    }
}
class Box57 {
    int value;
    int weight;
    Box57? next;

    function init(int v) {
        self.value = v + K57;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box57 n) {
        self.next = n;
    }
}
class Box58 {
    int value;
    int weight;
    Box58? next;

    function init(int v) {
        self.value = v + K58;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box58 n) {
        self.next = n;
    }
}
class Box59 {
    int value;
    int weight;
    Box59? next;

    function init(int v) {
        self.value = v + K59;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box59 n) {
        self.next = n;
    }
}
class Box60 {
    int value;
    int weight;
    Box60? next;

    function init(int v) {
        self.value = v + K60;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box60 n) {
        self.next = n;
    }
}
class Box61 {
    int value;
    int weight;
    Box61? next;

    function init(int v) {
        self.value = v + K61;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box61 n) {
        self.next = n;
    }
}
class Box62 {
    int value;
    int weight;
    Box62? next;

    function init(int v) {
        self.value = v + K62;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box62 n) {
        self.next = n;
    }
}
class Box63 {
    int value;
    int weight;
    Box63? next;

    function init(int v) {
        self.value = v + K63;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box63 n) {
        self.next = n;
    }
}
class Box64 {
    int value;
    int weight;
    Box64? next;

    function init(int v) {
        self.value = v + K64;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box64 n) {
        self.next = n;
    }
}
class Box65 {
    int value;
    int weight;
    Box65? next;

    function init(int v) {
        self.value = v + K65;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box65 n) {
        self.next = n;
    }
}
class Box66 {
    int value;
    int weight;
    Box66? next;

    function init(int v) {
        self.value = v + K66;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box66 n) {
        self.next = n;
    }
}
class Box67 {
    int value;
    int weight;
    Box67? next;

    function init(int v) {
        self.value = v + K67;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box67 n) {
        self.next = n;
    }
}
class Box68 {
    int value;
    int weight;
    Box68? next;

    function init(int v) {
        self.value = v + K68;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box68 n) {
        self.next = n;
    }
}
class Box69 {
    int value;
    int weight;
    Box69? next;

    function init(int v) {
        self.value = v + K69;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box69 n) {
        self.next = n;
    }
}
class Box70 {
    int value;
    int weight;
    Box70? next;

    function init(int v) {
        self.value = v + K70;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box70 n) {
        self.next = n;
    }
}
class Box71 {
    int value;
    int weight;
    Box71? next;

    function init(int v) {
        self.value = v + K71;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box71 n) {
        self.next = n;
    }
}
class Box72 {
    int value;
    int weight;
    Box72? next;

    function init(int v) {
        self.value = v + K72;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box72 n) {
        self.next = n;
    }
}
class Box73 {
    int value;
    int weight;
    Box73? next;

    function init(int v) {
        self.value = v + K73;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box73 n) {
        self.next = n;
    }
}
class Box74 {
    int value;
    int weight;
    Box74? next;

    function init(int v) {
        self.value = v + K74;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box74 n) {
        self.next = n;
    }
}
class Box75 {
    int value;
    int weight;
    Box75? next;

    function init(int v) {
        self.value = v + K75;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box75 n) {
        self.next = n;
    }
}
class Box76 {
    int value;
    int weight;
    Box76? next;

    function init(int v) {
        self.value = v + K76;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box76 n) {
        self.next = n;
    }
}
class Box77 {
    int value;
    int weight;
    Box77? next;

    function init(int v) {
        self.value = v + K77;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box77 n) {
        self.next = n;
    }
}
class Box78 {
    int value;
    int weight;
    Box78? next;

    function init(int v) {
        self.value = v + K78;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box78 n) {
        self.next = n;
    }
}
class Box79 {
    int value;
    int weight;
    Box79? next;

    function init(int v) {
        self.value = v + K79;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box79 n) {
        self.next = n;
    }
}
class Box80 {
    int value;
    int weight;
    Box80? next;

    function init(int v) {
        self.value = v + K80;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box80 n) {
        self.next = n;
    }
}
class Box81 {
    int value;
    int weight;
    Box81? next;

    function init(int v) {
        self.value = v + K81;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box81 n) {
        self.next = n;
    }
}
class Box82 {
    int value;
    int weight;
    Box82? next;

    function init(int v) {
        self.value = v + K82;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box82 n) {
        self.next = n;
    }
}
class Box83 {
    int value;
    int weight;
    Box83? next;

    function init(int v) {
        self.value = v + K83;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box83 n) {
        self.next = n;
    }
}
class Box84 {
    int value;
    int weight;
    Box84? next;

    function init(int v) {
        self.value = v + K84;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box84 n) {
        self.next = n;
    }
}
class Box85 {
    int value;
    int weight;
    Box85? next;

    function init(int v) {
        self.value = v + K85;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box85 n) {
        self.next = n;
    }
}
class Box86 {
    int value;
    int weight;
    Box86? next;

    function init(int v) {
        self.value = v + K86;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box86 n) {
        self.next = n;
    }
}
class Box87 {
    int value;
    int weight;
    Box87? next;

    function init(int v) {
        self.value = v + K87;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box87 n) {
        self.next = n;
    }
}
class Box88 {
    int value;
    int weight;
    Box88? next;

    function init(int v) {
        self.value = v + K88;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box88 n) {
        self.next = n;
    }
}
class Box89 {
    int value;
    int weight;
    Box89? next;

    function init(int v) {
        self.value = v + K89;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box89 n) {
        self.next = n;
    }
}
class Box90 {
    int value;
    int weight;
    Box90? next;

    function init(int v) {
        self.value = v + K90;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box90 n) {
        self.next = n;
    }
}
class Box91 {
    int value;
    int weight;
    Box91? next;

    function init(int v) {
        self.value = v + K91;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box91 n) {
        self.next = n;
    }
}
class Box92 {
    int value;
    int weight;
    Box92? next;

    function init(int v) {
        self.value = v + K92;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box92 n) {
        self.next = n;
    }
}
class Box93 {
    int value;
    int weight;
    Box93? next;

    function init(int v) {
        self.value = v + K93;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box93 n) {
        self.next = n;
    }
}
class Box94 {
    int value;
    int weight;
    Box94? next;

    function init(int v) {
        self.value = v + K94;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box94 n) {
        self.next = n;
    }
}
class Box95 {
    int value;
    int weight;
    Box95? next;

    function init(int v) {
        self.value = v + K95;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box95 n) {
        self.next = n;
    }
}
class Box96 {
    int value;
    int weight;
    Box96? next;

    function init(int v) {
        self.value = v + K96;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box96 n) {
        self.next = n;
    }
}
class Box97 {
    int value;
    int weight;
    Box97? next;

    function init(int v) {
        self.value = v + K97;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box97 n) {
        self.next = n;
    }
}
class Box98 {
    int value;
    int weight;
    Box98? next;

    function init(int v) {
        self.value = v + K98;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box98 n) {
        self.next = n;
    }
}
class Box99 {
    int value;
    int weight;
    Box99? next;

    function init(int v) {
        self.value = v + K99;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box99 n) {
        self.next = n;
    }
}
class Box100 {
    int value;
    int weight;
    Box100? next;

    function init(int v) {
        self.value = v + K100;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box100 n) {
        self.next = n;
    }
}
class Box101 {
    int value;
    int weight;
    Box101? next;

    function init(int v) {
        self.value = v + K101;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box101 n) {
        self.next = n;
    }
}
class Box102 {
    int value;
    int weight;
    Box102? next;

    function init(int v) {
        self.value = v + K102;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box102 n) {
        self.next = n;
    }
}
class Box103 {
    int value;
    int weight;
    Box103? next;

    function init(int v) {
        self.value = v + K103;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box103 n) {
        self.next = n;
    }
}
class Box104 {
    int value;
    int weight;
    Box104? next;

    function init(int v) {
        self.value = v + K104;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box104 n) {
        self.next = n;
    }
}
class Box105 {
    int value;
    int weight;
    Box105? next;

    function init(int v) {
        self.value = v + K105;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box105 n) {
        self.next = n;
    }
}
class Box106 {
    int value;
    int weight;
    Box106? next;

    function init(int v) {
        self.value = v + K106;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box106 n) {
        self.next = n;
    }
}
class Box107 {
    int value;
    int weight;
    Box107? next;

    function init(int v) {
        self.value = v + K107;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box107 n) {
        self.next = n;
    }
}
class Box108 {
    int value;
    int weight;
    Box108? next;

    function init(int v) {
        self.value = v + K108;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box108 n) {
        self.next = n;
    }
}
class Box109 {
    int value;
    int weight;
    Box109? next;

    function init(int v) {
        self.value = v + K109;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box109 n) {
        self.next = n;
    }
}
class Box110 {
    int value;
    int weight;
    Box110? next;

    function init(int v) {
        self.value = v + K110;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box110 n) {
        self.next = n;
    }
}
class Box111 {
    int value;
    int weight;
    Box111? next;

    function init(int v) {
        self.value = v + K111;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box111 n) {
        self.next = n;
    }
}
class Box112 {
    int value;
    int weight;
    Box112? next;

    function init(int v) {
        self.value = v + K112;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box112 n) {
        self.next = n;
    }
}
class Box113 {
    int value;
    int weight;
    Box113? next;

    function init(int v) {
        self.value = v + K113;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box113 n) {
        self.next = n;
    }
}
class Box114 {
    int value;
    int weight;
    Box114? next;

    function init(int v) {
        self.value = v + K114;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box114 n) {
        self.next = n;
    }
}
class Box115 {
    int value;
    int weight;
    Box115? next;

    function init(int v) {
        self.value = v + K115;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box115 n) {
        self.next = n;
    }
}
class Box116 {
    int value;
    int weight;
    Box116? next;

    function init(int v) {
        self.value = v + K116;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box116 n) {
        self.next = n;
    }
}
class Box117 {
    int value;
    int weight;
    Box117? next;

    function init(int v) {
        self.value = v + K117;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box117 n) {
        self.next = n;
    }
}
class Box118 {
    int value;
    int weight;
    Box118? next;

    function init(int v) {
        self.value = v + K118;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box118 n) {
        self.next = n;
    }
}
class Box119 {
    int value;
    int weight;
    Box119? next;

    function init(int v) {
        self.value = v + K119;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box119 n) {
        self.next = n;
    }
}
class Box120 {
    int value;
    int weight;
    Box120? next;

    function init(int v) {
        self.value = v + K120;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box120 n) {
        self.next = n;
    }
}
class Box121 {
    int value;
    int weight;
    Box121? next;

    function init(int v) {
        self.value = v + K121;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box121 n) {
        self.next = n;
    }
}
class Box122 {
    int value;
    int weight;
    Box122? next;

    function init(int v) {
        self.value = v + K122;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box122 n) {
        self.next = n;
    }
}
class Box123 {
    int value;
    int weight;
    Box123? next;

    function init(int v) {
        self.value = v + K123;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box123 n) {
        self.next = n;
    }
}
class Box124 {
    int value;
    int weight;
    Box124? next;

    function init(int v) {
        self.value = v + K124;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box124 n) {
        self.next = n;
    }
}
class Box125 {
    int value;
    int weight;
    Box125? next;

    function init(int v) {
        self.value = v + K125;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box125 n) {
        self.next = n;
    }
}
class Box126 {
    int value;
    int weight;
    Box126? next;

    function init(int v) {
        self.value = v + K126;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box126 n) {
        self.next = n;
    }
}
class Box127 {
    int value;
    int weight;
    Box127? next;

    function init(int v) {
        self.value = v + K127;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box127 n) {
        self.next = n;
    }
}
class Box128 {
    int value;
    int weight;
    Box128? next;

    function init(int v) {
        self.value = v + K128;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box128 n) {
        self.next = n;
    }
}
class Box129 {
    int value;
    int weight;
    Box129? next;

    function init(int v) {
        self.value = v + K129;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box129 n) {
        self.next = n;
    }
}
class Box130 {
    int value;
    int weight;
    Box130? next;

    function init(int v) {
        self.value = v + K130;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box130 n) {
        self.next = n;
    }
}
class Box131 {
    int value;
    int weight;
    Box131? next;

    function init(int v) {
        self.value = v + K131;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box131 n) {
        self.next = n;
    }
}
class Box132 {
    int value;
    int weight;
    Box132? next;

    function init(int v) {
        self.value = v + K132;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box132 n) {
        self.next = n;
    }
}
class Box133 {
    int value;
    int weight;
    Box133? next;

    function init(int v) {
        self.value = v + K133;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box133 n) {
        self.next = n;
    }
}
class Box134 {
    int value;
    int weight;
    Box134? next;

    function init(int v) {
        self.value = v + K134;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box134 n) {
        self.next = n;
    }
}
class Box135 {
    int value;
    int weight;
    Box135? next;

    function init(int v) {
        self.value = v + K135;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box135 n) {
        self.next = n;
    }
}
class Box136 {
    int value;
    int weight;
    Box136? next;

    function init(int v) {
        self.value = v + K136;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box136 n) {
        self.next = n;
    }
}
class Box137 {
    int value;
    int weight;
    Box137? next;

    function init(int v) {
        self.value = v + K137;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box137 n) {
        self.next = n;
    }
}
class Box138 {
    int value;
    int weight;
    Box138? next;

    function init(int v) {
        self.value = v + K138;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box138 n) {
        self.next = n;
    }
}
class Box139 {
    int value;
    int weight;
    Box139? next;

    function init(int v) {
        self.value = v + K139;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box139 n) {
        self.next = n;
    }
}
class Box140 {
    int value;
    int weight;
    Box140? next;

    function init(int v) {
        self.value = v + K140;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box140 n) {
        self.next = n;
    }
}
class Box141 {
    int value;
    int weight;
    Box141? next;

    function init(int v) {
        self.value = v + K141;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box141 n) {
        self.next = n;
    }
}
class Box142 {
    int value;
    int weight;
    Box142? next;

    function init(int v) {
        self.value = v + K142;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box142 n) {
        self.next = n;
    }
}
class Box143 {
    int value;
    int weight;
    Box143? next;

    function init(int v) {
        self.value = v + K143;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box143 n) {
        self.next = n;
    }
}
class Box144 {
    int value;
    int weight;
    Box144? next;

    function init(int v) {
        self.value = v + K144;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box144 n) {
        self.next = n;
    }
}
class Box145 {
    int value;
    int weight;
    Box145? next;

    function init(int v) {
        self.value = v + K145;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box145 n) {
        self.next = n;
    }
}
class Box146 {
    int value;
    int weight;
    Box146? next;

    function init(int v) {
        self.value = v + K146;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box146 n) {
        self.next = n;
    }
}
class Box147 {
    int value;
    int weight;
    Box147? next;

    function init(int v) {
        self.value = v + K147;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box147 n) {
        self.next = n;
    }
}
class Box148 {
    int value;
    int weight;
    Box148? next;

    function init(int v) {
        self.value = v + K148;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box148 n) {
        self.next = n;
    }
}
class Box149 {
    int value;
    int weight;
    Box149? next;

    function init(int v) {
        self.value = v + K149;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box149 n) {
        self.next = n;
    }
}
class Box150 {
    int value;
    int weight;
    Box150? next;

    function init(int v) {
        self.value = v + K150;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box150 n) {
        self.next = n;
    }
}
class Box151 {
    int value;
    int weight;
    Box151? next;

    function init(int v) {
        self.value = v + K151;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box151 n) {
        self.next = n;
    }
}
class Box152 {
    int value;
    int weight;
    Box152? next;

    function init(int v) {
        self.value = v + K152;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box152 n) {
        self.next = n;
    }
}
class Box153 {
    int value;
    int weight;
    Box153? next;

    function init(int v) {
        self.value = v + K153;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box153 n) {
        self.next = n;
    }
}
class Box154 {
    int value;
    int weight;
    Box154? next;

    function init(int v) {
        self.value = v + K154;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box154 n) {
        self.next = n;
    }
}
class Box155 {
    int value;
    int weight;
    Box155? next;

    function init(int v) {
        self.value = v + K155;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box155 n) {
        self.next = n;
    }
}
class Box156 {
    int value;
    int weight;
    Box156? next;

    function init(int v) {
        self.value = v + K156;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box156 n) {
        self.next = n;
    }
}
class Box157 {
    int value;
    int weight;
    Box157? next;

    function init(int v) {
        self.value = v + K157;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box157 n) {
        self.next = n;
    }
}
class Box158 {
    int value;
    int weight;
    Box158? next;

    function init(int v) {
        self.value = v + K158;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box158 n) {
        self.next = n;
    }
}
class Box159 {
    int value;
    int weight;
    Box159? next;

    function init(int v) {
        self.value = v + K159;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box159 n) {
        self.next = n;
    }
}
class Box160 {
    int value;
    int weight;
    Box160? next;

    function init(int v) {
        self.value = v + K160;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box160 n) {
        self.next = n;
    }
}
class Box161 {
    int value;
    int weight;
    Box161? next;

    function init(int v) {
        self.value = v + K161;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box161 n) {
        self.next = n;
    }
}
class Box162 {
    int value;
    int weight;
    Box162? next;

    function init(int v) {
        self.value = v + K162;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box162 n) {
        self.next = n;
    }
}
class Box163 {
    int value;
    int weight;
    Box163? next;

    function init(int v) {
        self.value = v + K163;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box163 n) {
        self.next = n;
    }
}
class Box164 {
    int value;
    int weight;
    Box164? next;

    function init(int v) {
        self.value = v + K164;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box164 n) {
        self.next = n;
    }
}
class Box165 {
    int value;
    int weight;
    Box165? next;

    function init(int v) {
        self.value = v + K165;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box165 n) {
        self.next = n;
    }
}
class Box166 {
    int value;
    int weight;
    Box166? next;

    function init(int v) {
        self.value = v + K166;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box166 n) {
        self.next = n;
    }
}
class Box167 {
    int value;
    int weight;
    Box167? next;

    function init(int v) {
        self.value = v + K167;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box167 n) {
        self.next = n;
    }
}
class Box168 {
    int value;
    int weight;
    Box168? next;

    function init(int v) {
        self.value = v + K168;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box168 n) {
        self.next = n;
    }
}
class Box169 {
    int value;
    int weight;
    Box169? next;

    function init(int v) {
        self.value = v + K169;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box169 n) {
        self.next = n;
    }
}
class Box170 {
    int value;
    int weight;
    Box170? next;

    function init(int v) {
        self.value = v + K170;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box170 n) {
        self.next = n;
    }
}
class Box171 {
    int value;
    int weight;
    Box171? next;

    function init(int v) {
        self.value = v + K171;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box171 n) {
        self.next = n;
    }
}
class Box172 {
    int value;
    int weight;
    Box172? next;

    function init(int v) {
        self.value = v + K172;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box172 n) {
        self.next = n;
    }
}
class Box173 {
    int value;
    int weight;
    Box173? next;

    function init(int v) {
        self.value = v + K173;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box173 n) {
        self.next = n;
    }
}
class Box174 {
    int value;
    int weight;
    Box174? next;

    function init(int v) {
        self.value = v + K174;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box174 n) {
        self.next = n;
    }
}
class Box175 {
    int value;
    int weight;
    Box175? next;

    function init(int v) {
        self.value = v + K175;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box175 n) {
        self.next = n;
    }
}
class Box176 {
    int value;
    int weight;
    Box176? next;

    function init(int v) {
        self.value = v + K176;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box176 n) {
        self.next = n;
    }
}
class Box177 {
    int value;
    int weight;
    Box177? next;

    function init(int v) {
        self.value = v + K177;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box177 n) {
        self.next = n;
    }
}
class Box178 {
    int value;
    int weight;
    Box178? next;

    function init(int v) {
        self.value = v + K178;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box178 n) {
        self.next = n;
    }
}
class Box179 {
    int value;
    int weight;
    Box179? next;

    function init(int v) {
        self.value = v + K179;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box179 n) {
        self.next = n;
    }
}
class Box180 {
    int value;
    int weight;
    Box180? next;

    function init(int v) {
        self.value = v + K180;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box180 n) {
        self.next = n;
    }
}
class Box181 {
    int value;
    int weight;
    Box181? next;

    function init(int v) {
        self.value = v + K181;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box181 n) {
        self.next = n;
    }
}
class Box182 {
    int value;
    int weight;
    Box182? next;

    function init(int v) {
        self.value = v + K182;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box182 n) {
        self.next = n;
    }
}
class Box183 {
    int value;
    int weight;
    Box183? next;

    function init(int v) {
        self.value = v + K183;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box183 n) {
        self.next = n;
    }
}
class Box184 {
    int value;
    int weight;
    Box184? next;

    function init(int v) {
        self.value = v + K184;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box184 n) {
        self.next = n;
    }
}
class Box185 {
    int value;
    int weight;
    Box185? next;

    function init(int v) {
        self.value = v + K185;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box185 n) {
        self.next = n;
    }
}
class Box186 {
    int value;
    int weight;
    Box186? next;

    function init(int v) {
        self.value = v + K186;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box186 n) {
        self.next = n;
    }
}
class Box187 {
    int value;
    int weight;
    Box187? next;

    function init(int v) {
        self.value = v + K187;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box187 n) {
        self.next = n;
    }
}
class Box188 {
    int value;
    int weight;
    Box188? next;

    function init(int v) {
        self.value = v + K188;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box188 n) {
        self.next = n;
    }
}
class Box189 {
    int value;
    int weight;
    Box189? next;

    function init(int v) {
        self.value = v + K189;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box189 n) {
        self.next = n;
    }
}
class Box190 {
    int value;
    int weight;
    Box190? next;

    function init(int v) {
        self.value = v + K190;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box190 n) {
        self.next = n;
    }
}
class Box191 {
    int value;
    int weight;
    Box191? next;

    function init(int v) {
        self.value = v + K191;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box191 n) {
        self.next = n;
    }
}
class Box192 {
    int value;
    int weight;
    Box192? next;

    function init(int v) {
        self.value = v + K192;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box192 n) {
        self.next = n;
    }
}
class Box193 {
    int value;
    int weight;
    Box193? next;

    function init(int v) {
        self.value = v + K193;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box193 n) {
        self.next = n;
    }
}
class Box194 {
    int value;
    int weight;
    Box194? next;

    function init(int v) {
        self.value = v + K194;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box194 n) {
        self.next = n;
    }
}
class Box195 {
    int value;
    int weight;
    Box195? next;

    function init(int v) {
        self.value = v + K195;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box195 n) {
        self.next = n;
    }
}
class Box196 {
    int value;
    int weight;
    Box196? next;

    function init(int v) {
        self.value = v + K196;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box196 n) {
        self.next = n;
    }
}
class Box197 {
    int value;
    int weight;
    Box197? next;

    function init(int v) {
        self.value = v + K197;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box197 n) {
        self.next = n;
    }
}
class Box198 {
    int value;
    int weight;
    Box198? next;

    function init(int v) {
        self.value = v + K198;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box198 n) {
        self.next = n;
    }
}
class Box199 {
    int value;
    int weight;
    Box199? next;

    function init(int v) {
        self.value = v + K199;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box199 n) {
        self.next = n;
    }
}
class Box200 {
    int value;
    int weight;
    Box200? next;

    function init(int v) {
        self.value = v + K0;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box200 n) {
        self.next = n;
    }
}
class Box201 {
    int value;
    int weight;
    Box201? next;

    function init(int v) {
        self.value = v + K1;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box201 n) {
        self.next = n;
    }
}
class Box202 {
    int value;
    int weight;
    Box202? next;

    function init(int v) {
        self.value = v + K2;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box202 n) {
        self.next = n;
    }
}
class Box203 {
    int value;
    int weight;
    Box203? next;

    function init(int v) {
        self.value = v + K3;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box203 n) {
        self.next = n;
    }
}
class Box204 {
    int value;
    int weight;
    Box204? next;

    function init(int v) {
        self.value = v + K4;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box204 n) {
        self.next = n;
    }
}
class Box205 {
    int value;
    int weight;
    Box205? next;

    function init(int v) {
        self.value = v + K5;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box205 n) {
        self.next = n;
    }
}
class Box206 {
    int value;
    int weight;
    Box206? next;

    function init(int v) {
        self.value = v + K6;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box206 n) {
        self.next = n;
    }
}
class Box207 {
    int value;
    int weight;
    Box207? next;

    function init(int v) {
        self.value = v + K7;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box207 n) {
        self.next = n;
    }
}
class Box208 {
    int value;
    int weight;
    Box208? next;

    function init(int v) {
        self.value = v + K8;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box208 n) {
        self.next = n;
    }
}
class Box209 {
    int value;
    int weight;
    Box209? next;

    function init(int v) {
        self.value = v + K9;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box209 n) {
        self.next = n;
    }
}
class Box210 {
    int value;
    int weight;
    Box210? next;

    function init(int v) {
        self.value = v + K10;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box210 n) {
        self.next = n;
    }
}
class Box211 {
    int value;
    int weight;
    Box211? next;

    function init(int v) {
        self.value = v + K11;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box211 n) {
        self.next = n;
    }
}
class Box212 {
    int value;
    int weight;
    Box212? next;

    function init(int v) {
        self.value = v + K12;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box212 n) {
        self.next = n;
    }
}
class Box213 {
    int value;
    int weight;
    Box213? next;

    function init(int v) {
        self.value = v + K13;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box213 n) {
        self.next = n;
    }
}
class Box214 {
    int value;
    int weight;
    Box214? next;

    function init(int v) {
        self.value = v + K14;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box214 n) {
        self.next = n;
    }
}
class Box215 {
    int value;
    int weight;
    Box215? next;

    function init(int v) {
        self.value = v + K15;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box215 n) {
        self.next = n;
    }
}
class Box216 {
    int value;
    int weight;
    Box216? next;

    function init(int v) {
        self.value = v + K16;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box216 n) {
        self.next = n;
    }
}
class Box217 {
    int value;
    int weight;
    Box217? next;

    function init(int v) {
        self.value = v + K17;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box217 n) {
        self.next = n;
    }
}
class Box218 {
    int value;
    int weight;
    Box218? next;

    function init(int v) {
        self.value = v + K18;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box218 n) {
        self.next = n;
    }
}
class Box219 {
    int value;
    int weight;
    Box219? next;

    function init(int v) {
        self.value = v + K19;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box219 n) {
        self.next = n;
    }
}
class Box220 {
    int value;
    int weight;
    Box220? next;

    function init(int v) {
        self.value = v + K20;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box220 n) {
        self.next = n;
    }
}
class Box221 {
    int value;
    int weight;
    Box221? next;

    function init(int v) {
        self.value = v + K21;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box221 n) {
        self.next = n;
    }
}
class Box222 {
    int value;
    int weight;
    Box222? next;

    function init(int v) {
        self.value = v + K22;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box222 n) {
        self.next = n;
    }
}
class Box223 {
    int value;
    int weight;
    Box223? next;

    function init(int v) {
        self.value = v + K23;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box223 n) {
        self.next = n;
    }
}
class Box224 {
    int value;
    int weight;
    Box224? next;

    function init(int v) {
        self.value = v + K24;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box224 n) {
        self.next = n;
    }
}
class Box225 {
    int value;
    int weight;
    Box225? next;

    function init(int v) {
        self.value = v + K25;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box225 n) {
        self.next = n;
    }
}
class Box226 {
    int value;
    int weight;
    Box226? next;

    function init(int v) {
        self.value = v + K26;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box226 n) {
        self.next = n;
    }
}
class Box227 {
    int value;
    int weight;
    Box227? next;

    function init(int v) {
        self.value = v + K27;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box227 n) {
        self.next = n;
    }
}
class Box228 {
    int value;
    int weight;
    Box228? next;

    function init(int v) {
        self.value = v + K28;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box228 n) {
        self.next = n;
    }
}
class Box229 {
    int value;
    int weight;
    Box229? next;

    function init(int v) {
        self.value = v + K29;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box229 n) {
        self.next = n;
    }
}
class Box230 {
    int value;
    int weight;
    Box230? next;

    function init(int v) {
        self.value = v + K30;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box230 n) {
        self.next = n;
    }
}
class Box231 {
    int value;
    int weight;
    Box231? next;

    function init(int v) {
        self.value = v + K31;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box231 n) {
        self.next = n;
    }
}
class Box232 {
    int value;
    int weight;
    Box232? next;

    function init(int v) {
        self.value = v + K32;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box232 n) {
        self.next = n;
    }
}
class Box233 {
    int value;
    int weight;
    Box233? next;

    function init(int v) {
        self.value = v + K33;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box233 n) {
        self.next = n;
    }
}
class Box234 {
    int value;
    int weight;
    Box234? next;

    function init(int v) {
        self.value = v + K34;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box234 n) {
        self.next = n;
    }
}
class Box235 {
    int value;
    int weight;
    Box235? next;

    function init(int v) {
        self.value = v + K35;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box235 n) {
        self.next = n;
    }
}
class Box236 {
    int value;
    int weight;
    Box236? next;

    function init(int v) {
        self.value = v + K36;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box236 n) {
        self.next = n;
    }
}
class Box237 {
    int value;
    int weight;
    Box237? next;

    function init(int v) {
        self.value = v + K37;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box237 n) {
        self.next = n;
    }
}
class Box238 {
    int value;
    int weight;
    Box238? next;

    function init(int v) {
        self.value = v + K38;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box238 n) {
        self.next = n;
    }
}
class Box239 {
    int value;
    int weight;
    Box239? next;

    function init(int v) {
        self.value = v + K39;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box239 n) {
        self.next = n;
    }
}
class Box240 {
    int value;
    int weight;
    Box240? next;

    function init(int v) {
        self.value = v + K40;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box240 n) {
        self.next = n;
    }
}
class Box241 {
    int value;
    int weight;
    Box241? next;

    function init(int v) {
        self.value = v + K41;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box241 n) {
        self.next = n;
    }
}
class Box242 {
    int value;
    int weight;
    Box242? next;

    function init(int v) {
        self.value = v + K42;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box242 n) {
        self.next = n;
    }
}
class Box243 {
    int value;
    int weight;
    Box243? next;

    function init(int v) {
        self.value = v + K43;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box243 n) {
        self.next = n;
    }
}
class Box244 {
    int value;
    int weight;
    Box244? next;

    function init(int v) {
        self.value = v + K44;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box244 n) {
        self.next = n;
    }
}
class Box245 {
    int value;
    int weight;
    Box245? next;

    function init(int v) {
        self.value = v + K45;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box245 n) {
        self.next = n;
    }
}
class Box246 {
    int value;
    int weight;
    Box246? next;

    function init(int v) {
        self.value = v + K46;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box246 n) {
        self.next = n;
    }
}
class Box247 {
    int value;
    int weight;
    Box247? next;

    function init(int v) {
        self.value = v + K47;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box247 n) {
        self.next = n;
    }
}
class Box248 {
    int value;
    int weight;
    Box248? next;

    function init(int v) {
        self.value = v + K48;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box248 n) {
        self.next = n;
    }
}
class Box249 {
    int value;
    int weight;
    Box249? next;

    function init(int v) {
        self.value = v + K49;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box249 n) {
        self.next = n;
    }
}
class Box250 {
    int value;
    int weight;
    Box250? next;

    function init(int v) {
        self.value = v + K50;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box250 n) {
        self.next = n;
    }
}
class Box251 {
    int value;
    int weight;
    Box251? next;

    function init(int v) {
        self.value = v + K51;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box251 n) {
        self.next = n;
    }
}
class Box252 {
    int value;
    int weight;
    Box252? next;

    function init(int v) {
        self.value = v + K52;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box252 n) {
        self.next = n;
    }
}
class Box253 {
    int value;
    int weight;
    Box253? next;

    function init(int v) {
        self.value = v + K53;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box253 n) {
        self.next = n;
    }
}
class Box254 {
    int value;
    int weight;
    Box254? next;

    function init(int v) {
        self.value = v + K54;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box254 n) {
        self.next = n;
    }
}
class Box255 {
    int value;
    int weight;
    Box255? next;

    function init(int v) {
        self.value = v + K55;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box255 n) {
        self.next = n;
    }
}
class Box256 {
    int value;
    int weight;
    Box256? next;

    function init(int v) {
        self.value = v + K56;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box256 n) {
        self.next = n;
    }
}
class Box257 {
    int value;
    int weight;
    Box257? next;

    function init(int v) {
        self.value = v + K57;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box257 n) {
        self.next = n;
    }
}
class Box258 {
    int value;
    int weight;
    Box258? next;

    function init(int v) {
        self.value = v + K58;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box258 n) {
        self.next = n;
    }
}
class Box259 {
    int value;
    int weight;
    Box259? next;

    function init(int v) {
        self.value = v + K59;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box259 n) {
        self.next = n;
    }
}
class Box260 {
    int value;
    int weight;
    Box260? next;

    function init(int v) {
        self.value = v + K60;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box260 n) {
        self.next = n;
    }
}
class Box261 {
    int value;
    int weight;
    Box261? next;

    function init(int v) {
        self.value = v + K61;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box261 n) {
        self.next = n;
    }
}
class Box262 {
    int value;
    int weight;
    Box262? next;

    function init(int v) {
        self.value = v + K62;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box262 n) {
        self.next = n;
    }
}
class Box263 {
    int value;
    int weight;
    Box263? next;

    function init(int v) {
        self.value = v + K63;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box263 n) {
        self.next = n;
    }
}
class Box264 {
    int value;
    int weight;
    Box264? next;

    function init(int v) {
        self.value = v + K64;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box264 n) {
        self.next = n;
    }
}
class Box265 {
    int value;
    int weight;
    Box265? next;

    function init(int v) {
        self.value = v + K65;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box265 n) {
        self.next = n;
    }
}
class Box266 {
    int value;
    int weight;
    Box266? next;

    function init(int v) {
        self.value = v + K66;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box266 n) {
        self.next = n;
    }
}
class Box267 {
    int value;
    int weight;
    Box267? next;

    function init(int v) {
        self.value = v + K67;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box267 n) {
        self.next = n;
    }
}
class Box268 {
    int value;
    int weight;
    Box268? next;

    function init(int v) {
        self.value = v + K68;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box268 n) {
        self.next = n;
    }
}
class Box269 {
    int value;
    int weight;
    Box269? next;

    function init(int v) {
        self.value = v + K69;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box269 n) {
        self.next = n;
    }
}
class Box270 {
    int value;
    int weight;
    Box270? next;

    function init(int v) {
        self.value = v + K70;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box270 n) {
        self.next = n;
    }
}
class Box271 {
    int value;
    int weight;
    Box271? next;

    function init(int v) {
        self.value = v + K71;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box271 n) {
        self.next = n;
    }
}
class Box272 {
    int value;
    int weight;
    Box272? next;

    function init(int v) {
        self.value = v + K72;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box272 n) {
        self.next = n;
    }
}
class Box273 {
    int value;
    int weight;
    Box273? next;

    function init(int v) {
        self.value = v + K73;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box273 n) {
        self.next = n;
    }
}
class Box274 {
    int value;
    int weight;
    Box274? next;

    function init(int v) {
        self.value = v + K74;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box274 n) {
        self.next = n;
    }
}
class Box275 {
    int value;
    int weight;
    Box275? next;

    function init(int v) {
        self.value = v + K75;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box275 n) {
        self.next = n;
    }
}
class Box276 {
    int value;
    int weight;
    Box276? next;

    function init(int v) {
        self.value = v + K76;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box276 n) {
        self.next = n;
    }
}
class Box277 {
    int value;
    int weight;
    Box277? next;

    function init(int v) {
        self.value = v + K77;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box277 n) {
        self.next = n;
    }
}
class Box278 {
    int value;
    int weight;
    Box278? next;

    function init(int v) {
        self.value = v + K78;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box278 n) {
        self.next = n;
    }
}
class Box279 {
    int value;
    int weight;
    Box279? next;

    function init(int v) {
        self.value = v + K79;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box279 n) {
        self.next = n;
    }
}
class Box280 {
    int value;
    int weight;
    Box280? next;

    function init(int v) {
        self.value = v + K80;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box280 n) {
        self.next = n;
    }
}
class Box281 {
    int value;
    int weight;
    Box281? next;

    function init(int v) {
        self.value = v + K81;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box281 n) {
        self.next = n;
    }
}
class Box282 {
    int value;
    int weight;
    Box282? next;

    function init(int v) {
        self.value = v + K82;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box282 n) {
        self.next = n;
    }
}
class Box283 {
    int value;
    int weight;
    Box283? next;

    function init(int v) {
        self.value = v + K83;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box283 n) {
        self.next = n;
    }
}
class Box284 {
    int value;
    int weight;
    Box284? next;

    function init(int v) {
        self.value = v + K84;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box284 n) {
        self.next = n;
    }
}
class Box285 {
    int value;
    int weight;
    Box285? next;

    function init(int v) {
        self.value = v + K85;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box285 n) {
        self.next = n;
    }
}
class Box286 {
    int value;
    int weight;
    Box286? next;

    function init(int v) {
        self.value = v + K86;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box286 n) {
        self.next = n;
    }
}
class Box287 {
    int value;
    int weight;
    Box287? next;

    function init(int v) {
        self.value = v + K87;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box287 n) {
        self.next = n;
    }
}
class Box288 {
    int value;
    int weight;
    Box288? next;

    function init(int v) {
        self.value = v + K88;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box288 n) {
        self.next = n;
    }
}
class Box289 {
    int value;
    int weight;
    Box289? next;

    function init(int v) {
        self.value = v + K89;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box289 n) {
        self.next = n;
    }
}
class Box290 {
    int value;
    int weight;
    Box290? next;

    function init(int v) {
        self.value = v + K90;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box290 n) {
        self.next = n;
    }
}
class Box291 {
    int value;
    int weight;
    Box291? next;

    function init(int v) {
        self.value = v + K91;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box291 n) {
        self.next = n;
    }
}
class Box292 {
    int value;
    int weight;
    Box292? next;

    function init(int v) {
        self.value = v + K92;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box292 n) {
        self.next = n;
    }
}
class Box293 {
    int value;
    int weight;
    Box293? next;

    function init(int v) {
        self.value = v + K93;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box293 n) {
        self.next = n;
    }
}
class Box294 {
    int value;
    int weight;
    Box294? next;

    function init(int v) {
        self.value = v + K94;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box294 n) {
        self.next = n;
    }
}
class Box295 {
    int value;
    int weight;
    Box295? next;

    function init(int v) {
        self.value = v + K95;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box295 n) {
        self.next = n;
    }
}
class Box296 {
    int value;
    int weight;
    Box296? next;

    function init(int v) {
        self.value = v + K96;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box296 n) {
        self.next = n;
    }
}
class Box297 {
    int value;
    int weight;
    Box297? next;

    function init(int v) {
        self.value = v + K97;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box297 n) {
        self.next = n;
    }
}
class Box298 {
    int value;
    int weight;
    Box298? next;

    function init(int v) {
        self.value = v + K98;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box298 n) {
        self.next = n;
    }
}
class Box299 {
    int value;
    int weight;
    Box299? next;

    function init(int v) {
        self.value = v + K99;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box299 n) {
        self.next = n;
    }
}
class Box300 {
    int value;
    int weight;
    Box300? next;

    function init(int v) {
        self.value = v + K100;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box300 n) {
        self.next = n;
    }
}
class Box301 {
    int value;
    int weight;
    Box301? next;

    function init(int v) {
        self.value = v + K101;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box301 n) {
        self.next = n;
    }
}
class Box302 {
    int value;
    int weight;
    Box302? next;

    function init(int v) {
        self.value = v + K102;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box302 n) {
        self.next = n;
    }
}
class Box303 {
    int value;
    int weight;
    Box303? next;

    function init(int v) {
        self.value = v + K103;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box303 n) {
        self.next = n;
    }
}
class Box304 {
    int value;
    int weight;
    Box304? next;

    function init(int v) {
        self.value = v + K104;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box304 n) {
        self.next = n;
    }
}
class Box305 {
    int value;
    int weight;
    Box305? next;

    function init(int v) {
        self.value = v + K105;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box305 n) {
        self.next = n;
    }
}
class Box306 {
    int value;
    int weight;
    Box306? next;

    function init(int v) {
        self.value = v + K106;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box306 n) {
        self.next = n;
    }
}
class Box307 {
    int value;
    int weight;
    Box307? next;

    function init(int v) {
        self.value = v + K107;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box307 n) {
        self.next = n;
    }
}
class Box308 {
    int value;
    int weight;
    Box308? next;

    function init(int v) {
        self.value = v + K108;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box308 n) {
        self.next = n;
    }
}
class Box309 {
    int value;
    int weight;
    Box309? next;

    function init(int v) {
        self.value = v + K109;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box309 n) {
        self.next = n;
    }
}
class Box310 {
    int value;
    int weight;
    Box310? next;

    function init(int v) {
        self.value = v + K110;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box310 n) {
        self.next = n;
    }
}
class Box311 {
    int value;
    int weight;
    Box311? next;

    function init(int v) {
        self.value = v + K111;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box311 n) {
        self.next = n;
    }
}
class Box312 {
    int value;
    int weight;
    Box312? next;

    function init(int v) {
        self.value = v + K112;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box312 n) {
        self.next = n;
    }
}
class Box313 {
    int value;
    int weight;
    Box313? next;

    function init(int v) {
        self.value = v + K113;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box313 n) {
        self.next = n;
    }
}
class Box314 {
    int value;
    int weight;
    Box314? next;

    function init(int v) {
        self.value = v + K114;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box314 n) {
        self.next = n;
    }
}
class Box315 {
    int value;
    int weight;
    Box315? next;

    function init(int v) {
        self.value = v + K115;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box315 n) {
        self.next = n;
    }
}
class Box316 {
    int value;
    int weight;
    Box316? next;

    function init(int v) {
        self.value = v + K116;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box316 n) {
        self.next = n;
    }
}
class Box317 {
    int value;
    int weight;
    Box317? next;

    function init(int v) {
        self.value = v + K117;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box317 n) {
        self.next = n;
    }
}
class Box318 {
    int value;
    int weight;
    Box318? next;

    function init(int v) {
        self.value = v + K118;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box318 n) {
        self.next = n;
    }
}
class Box319 {
    int value;
    int weight;
    Box319? next;

    function init(int v) {
        self.value = v + K119;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box319 n) {
        self.next = n;
    }
}
class Box320 {
    int value;
    int weight;
    Box320? next;

    function init(int v) {
        self.value = v + K120;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box320 n) {
        self.next = n;
    }
}
class Box321 {
    int value;
    int weight;
    Box321? next;

    function init(int v) {
        self.value = v + K121;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box321 n) {
        self.next = n;
    }
}
class Box322 {
    int value;
    int weight;
    Box322? next;

    function init(int v) {
        self.value = v + K122;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box322 n) {
        self.next = n;
    }
}
class Box323 {
    int value;
    int weight;
    Box323? next;

    function init(int v) {
        self.value = v + K123;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box323 n) {
        self.next = n;
    }
}
class Box324 {
    int value;
    int weight;
    Box324? next;

    function init(int v) {
        self.value = v + K124;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box324 n) {
        self.next = n;
    }
}
class Box325 {
    int value;
    int weight;
    Box325? next;

    function init(int v) {
        self.value = v + K125;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box325 n) {
        self.next = n;
    }
}
class Box326 {
    int value;
    int weight;
    Box326? next;

    function init(int v) {
        self.value = v + K126;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box326 n) {
        self.next = n;
    }
}
class Box327 {
    int value;
    int weight;
    Box327? next;

    function init(int v) {
        self.value = v + K127;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box327 n) {
        self.next = n;
    }
}
class Box328 {
    int value;
    int weight;
    Box328? next;

    function init(int v) {
        self.value = v + K128;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box328 n) {
        self.next = n;
    }
}
class Box329 {
    int value;
    int weight;
    Box329? next;

    function init(int v) {
        self.value = v + K129;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box329 n) {
        self.next = n;
    }
}
class Box330 {
    int value;
    int weight;
    Box330? next;

    function init(int v) {
        self.value = v + K130;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box330 n) {
        self.next = n;
    }
}
class Box331 {
    int value;
    int weight;
    Box331? next;

    function init(int v) {
        self.value = v + K131;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box331 n) {
        self.next = n;
    }
}
class Box332 {
    int value;
    int weight;
    Box332? next;

    function init(int v) {
        self.value = v + K132;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box332 n) {
        self.next = n;
    }
}
class Box333 {
    int value;
    int weight;
    Box333? next;

    function init(int v) {
        self.value = v + K133;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box333 n) {
        self.next = n;
    }
}
class Box334 {
    int value;
    int weight;
    Box334? next;

    function init(int v) {
        self.value = v + K134;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box334 n) {
        self.next = n;
    }
}
class Box335 {
    int value;
    int weight;
    Box335? next;

    function init(int v) {
        self.value = v + K135;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box335 n) {
        self.next = n;
    }
}
class Box336 {
    int value;
    int weight;
    Box336? next;

    function init(int v) {
        self.value = v + K136;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box336 n) {
        self.next = n;
    }
}
class Box337 {
    int value;
    int weight;
    Box337? next;

    function init(int v) {
        self.value = v + K137;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box337 n) {
        self.next = n;
    }
}
class Box338 {
    int value;
    int weight;
    Box338? next;

    function init(int v) {
        self.value = v + K138;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box338 n) {
        self.next = n;
    }
}
class Box339 {
    int value;
    int weight;
    Box339? next;

    function init(int v) {
        self.value = v + K139;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box339 n) {
        self.next = n;
    }
}
class Box340 {
    int value;
    int weight;
    Box340? next;

    function init(int v) {
        self.value = v + K140;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box340 n) {
        self.next = n;
    }
}
class Box341 {
    int value;
    int weight;
    Box341? next;

    function init(int v) {
        self.value = v + K141;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box341 n) {
        self.next = n;
    }
}
class Box342 {
    int value;
    int weight;
    Box342? next;

    function init(int v) {
        self.value = v + K142;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box342 n) {
        self.next = n;
    }
}
class Box343 {
    int value;
    int weight;
    Box343? next;

    function init(int v) {
        self.value = v + K143;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box343 n) {
        self.next = n;
    }
}
class Box344 {
    int value;
    int weight;
    Box344? next;

    function init(int v) {
        self.value = v + K144;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box344 n) {
        self.next = n;
    }
}
class Box345 {
    int value;
    int weight;
    Box345? next;

    function init(int v) {
        self.value = v + K145;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box345 n) {
        self.next = n;
    }
}
class Box346 {
    int value;
    int weight;
    Box346? next;

    function init(int v) {
        self.value = v + K146;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box346 n) {
        self.next = n;
    }
}
class Box347 {
    int value;
    int weight;
    Box347? next;

    function init(int v) {
        self.value = v + K147;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box347 n) {
        self.next = n;
    }
}
class Box348 {
    int value;
    int weight;
    Box348? next;

    function init(int v) {
        self.value = v + K148;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box348 n) {
        self.next = n;
    }
}
class Box349 {
    int value;
    int weight;
    Box349? next;

    function init(int v) {
        self.value = v + K149;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box349 n) {
        self.next = n;
    }
}
class Box350 {
    int value;
    int weight;
    Box350? next;

    function init(int v) {
        self.value = v + K150;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box350 n) {
        self.next = n;
    }
}
class Box351 {
    int value;
    int weight;
    Box351? next;

    function init(int v) {
        self.value = v + K151;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box351 n) {
        self.next = n;
    }
}
class Box352 {
    int value;
    int weight;
    Box352? next;

    function init(int v) {
        self.value = v + K152;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box352 n) {
        self.next = n;
    }
}
class Box353 {
    int value;
    int weight;
    Box353? next;

    function init(int v) {
        self.value = v + K153;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box353 n) {
        self.next = n;
    }
}
class Box354 {
    int value;
    int weight;
    Box354? next;

    function init(int v) {
        self.value = v + K154;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box354 n) {
        self.next = n;
    }
}
class Box355 {
    int value;
    int weight;
    Box355? next;

    function init(int v) {
        self.value = v + K155;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box355 n) {
        self.next = n;
    }
}
class Box356 {
    int value;
    int weight;
    Box356? next;

    function init(int v) {
        self.value = v + K156;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box356 n) {
        self.next = n;
    }
}
class Box357 {
    int value;
    int weight;
    Box357? next;

    function init(int v) {
        self.value = v + K157;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box357 n) {
        self.next = n;
    }
}
class Box358 {
    int value;
    int weight;
    Box358? next;

    function init(int v) {
        self.value = v + K158;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box358 n) {
        self.next = n;
    }
}
class Box359 {
    int value;
    int weight;
    Box359? next;

    function init(int v) {
        self.value = v + K159;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box359 n) {
        self.next = n;
    }
}
class Box360 {
    int value;
    int weight;
    Box360? next;

    function init(int v) {
        self.value = v + K160;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box360 n) {
        self.next = n;
    }
}
class Box361 {
    int value;
    int weight;
    Box361? next;

    function init(int v) {
        self.value = v + K161;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box361 n) {
        self.next = n;
    }
}
class Box362 {
    int value;
    int weight;
    Box362? next;

    function init(int v) {
        self.value = v + K162;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box362 n) {
        self.next = n;
    }
}
class Box363 {
    int value;
    int weight;
    Box363? next;

    function init(int v) {
        self.value = v + K163;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box363 n) {
        self.next = n;
    }
}
class Box364 {
    int value;
    int weight;
    Box364? next;

    function init(int v) {
        self.value = v + K164;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box364 n) {
        self.next = n;
    }
}
class Box365 {
    int value;
    int weight;
    Box365? next;

    function init(int v) {
        self.value = v + K165;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box365 n) {
        self.next = n;
    }
}
class Box366 {
    int value;
    int weight;
    Box366? next;

    function init(int v) {
        self.value = v + K166;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box366 n) {
        self.next = n;
    }
}
class Box367 {
    int value;
    int weight;
    Box367? next;

    function init(int v) {
        self.value = v + K167;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box367 n) {
        self.next = n;
    }
}
class Box368 {
    int value;
    int weight;
    Box368? next;

    function init(int v) {
        self.value = v + K168;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box368 n) {
        self.next = n;
    }
}
class Box369 {
    int value;
    int weight;
    Box369? next;

    function init(int v) {
        self.value = v + K169;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box369 n) {
        self.next = n;
    }
}
class Box370 {
    int value;
    int weight;
    Box370? next;

    function init(int v) {
        self.value = v + K170;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box370 n) {
        self.next = n;
    }
}
class Box371 {
    int value;
    int weight;
    Box371? next;

    function init(int v) {
        self.value = v + K171;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box371 n) {
        self.next = n;
    }
}
class Box372 {
    int value;
    int weight;
    Box372? next;

    function init(int v) {
        self.value = v + K172;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box372 n) {
        self.next = n;
    }
}
class Box373 {
    int value;
    int weight;
    Box373? next;

    function init(int v) {
        self.value = v + K173;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box373 n) {
        self.next = n;
    }
}
class Box374 {
    int value;
    int weight;
    Box374? next;

    function init(int v) {
        self.value = v + K174;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box374 n) {
        self.next = n;
    }
}
class Box375 {
    int value;
    int weight;
    Box375? next;

    function init(int v) {
        self.value = v + K175;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box375 n) {
        self.next = n;
    }
}
class Box376 {
    int value;
    int weight;
    Box376? next;

    function init(int v) {
        self.value = v + K176;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box376 n) {
        self.next = n;
    }
}
class Box377 {
    int value;
    int weight;
    Box377? next;

    function init(int v) {
        self.value = v + K177;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box377 n) {
        self.next = n;
    }
}
class Box378 {
    int value;
    int weight;
    Box378? next;

    function init(int v) {
        self.value = v + K178;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box378 n) {
        self.next = n;
    }
}
class Box379 {
    int value;
    int weight;
    Box379? next;

    function init(int v) {
        self.value = v + K179;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box379 n) {
        self.next = n;
    }
}
class Box380 {
    int value;
    int weight;
    Box380? next;

    function init(int v) {
        self.value = v + K180;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box380 n) {
        self.next = n;
    }
}
class Box381 {
    int value;
    int weight;
    Box381? next;

    function init(int v) {
        self.value = v + K181;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box381 n) {
        self.next = n;
    }
}
class Box382 {
    int value;
    int weight;
    Box382? next;

    function init(int v) {
        self.value = v + K182;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box382 n) {
        self.next = n;
    }
}
class Box383 {
    int value;
    int weight;
    Box383? next;

    function init(int v) {
        self.value = v + K183;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box383 n) {
        self.next = n;
    }
}
class Box384 {
    int value;
    int weight;
    Box384? next;

    function init(int v) {
        self.value = v + K184;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box384 n) {
        self.next = n;
    }
}
class Box385 {
    int value;
    int weight;
    Box385? next;

    function init(int v) {
        self.value = v + K185;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box385 n) {
        self.next = n;
    }
}
class Box386 {
    int value;
    int weight;
    Box386? next;

    function init(int v) {
        self.value = v + K186;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box386 n) {
        self.next = n;
    }
}
class Box387 {
    int value;
    int weight;
    Box387? next;

    function init(int v) {
        self.value = v + K187;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box387 n) {
        self.next = n;
    }
}
class Box388 {
    int value;
    int weight;
    Box388? next;

    function init(int v) {
        self.value = v + K188;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box388 n) {
        self.next = n;
    }
}
class Box389 {
    int value;
    int weight;
    Box389? next;

    function init(int v) {
        self.value = v + K189;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box389 n) {
        self.next = n;
    }
}
class Box390 {
    int value;
    int weight;
    Box390? next;

    function init(int v) {
        self.value = v + K190;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box390 n) {
        self.next = n;
    }
}
class Box391 {
    int value;
    int weight;
    Box391? next;

    function init(int v) {
        self.value = v + K191;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box391 n) {
        self.next = n;
    }
}
class Box392 {
    int value;
    int weight;
    Box392? next;

    function init(int v) {
        self.value = v + K192;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box392 n) {
        self.next = n;
    }
}
class Box393 {
    int value;
    int weight;
    Box393? next;

    function init(int v) {
        self.value = v + K193;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box393 n) {
        self.next = n;
    }
}
class Box394 {
    int value;
    int weight;
    Box394? next;

    function init(int v) {
        self.value = v + K194;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box394 n) {
        self.next = n;
    }
}
class Box395 {
    int value;
    int weight;
    Box395? next;

    function init(int v) {
        self.value = v + K195;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box395 n) {
        self.next = n;
    }
}
class Box396 {
    int value;
    int weight;
    Box396? next;

    function init(int v) {
        self.value = v + K196;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box396 n) {
        self.next = n;
    }
}
class Box397 {
    int value;
    int weight;
    Box397? next;

    function init(int v) {
        self.value = v + K197;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box397 n) {
        self.next = n;
    }
}
class Box398 {
    int value;
    int weight;
    Box398? next;

    function init(int v) {
        self.value = v + K198;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box398 n) {
        self.next = n;
    }
}
class Box399 {
    int value;
    int weight;
    Box399? next;

    function init(int v) {
        self.value = v + K199;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box399 n) {
        self.next = n;
    }
}
class Box400 {
    int value;
    int weight;
    Box400? next;

    function init(int v) {
        self.value = v + K0;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box400 n) {
        self.next = n;
    }
}
class Box401 {
    int value;
    int weight;
    Box401? next;

    function init(int v) {
        self.value = v + K1;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box401 n) {
        self.next = n;
    }
}
class Box402 {
    int value;
    int weight;
    Box402? next;

    function init(int v) {
        self.value = v + K2;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box402 n) {
        self.next = n;
    }
}
class Box403 {
    int value;
    int weight;
    Box403? next;

    function init(int v) {
        self.value = v + K3;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box403 n) {
        self.next = n;
    }
}
class Box404 {
    int value;
    int weight;
    Box404? next;

    function init(int v) {
        self.value = v + K4;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box404 n) {
        self.next = n;
    }
}
class Box405 {
    int value;
    int weight;
    Box405? next;

    function init(int v) {
        self.value = v + K5;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box405 n) {
        self.next = n;
    }
}
class Box406 {
    int value;
    int weight;
    Box406? next;

    function init(int v) {
        self.value = v + K6;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box406 n) {
        self.next = n;
    }
}
class Box407 {
    int value;
    int weight;
    Box407? next;

    function init(int v) {
        self.value = v + K7;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box407 n) {
        self.next = n;
    }
}
class Box408 {
    int value;
    int weight;
    Box408? next;

    function init(int v) {
        self.value = v + K8;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box408 n) {
        self.next = n;
    }
}
class Box409 {
    int value;
    int weight;
    Box409? next;

    function init(int v) {
        self.value = v + K9;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box409 n) {
        self.next = n;
    }
}
class Box410 {
    int value;
    int weight;
    Box410? next;

    function init(int v) {
        self.value = v + K10;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box410 n) {
        self.next = n;
    }
}
class Box411 {
    int value;
    int weight;
    Box411? next;

    function init(int v) {
        self.value = v + K11;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box411 n) {
        self.next = n;
    }
}
class Box412 {
    int value;
    int weight;
    Box412? next;

    function init(int v) {
        self.value = v + K12;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box412 n) {
        self.next = n;
    }
}
class Box413 {
    int value;
    int weight;
    Box413? next;

    function init(int v) {
        self.value = v + K13;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box413 n) {
        self.next = n;
    }
}
class Box414 {
    int value;
    int weight;
    Box414? next;

    function init(int v) {
        self.value = v + K14;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box414 n) {
        self.next = n;
    }
}
class Box415 {
    int value;
    int weight;
    Box415? next;

    function init(int v) {
        self.value = v + K15;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box415 n) {
        self.next = n;
    }
}
class Box416 {
    int value;
    int weight;
    Box416? next;

    function init(int v) {
        self.value = v + K16;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box416 n) {
        self.next = n;
    }
}
class Box417 {
    int value;
    int weight;
    Box417? next;

    function init(int v) {
        self.value = v + K17;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box417 n) {
        self.next = n;
    }
}
class Box418 {
    int value;
    int weight;
    Box418? next;

    function init(int v) {
        self.value = v + K18;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box418 n) {
        self.next = n;
    }
}
class Box419 {
    int value;
    int weight;
    Box419? next;

    function init(int v) {
        self.value = v + K19;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box419 n) {
        self.next = n;
    }
}
class Box420 {
    int value;
    int weight;
    Box420? next;

    function init(int v) {
        self.value = v + K20;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box420 n) {
        self.next = n;
    }
}
class Box421 {
    int value;
    int weight;
    Box421? next;

    function init(int v) {
        self.value = v + K21;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box421 n) {
        self.next = n;
    }
}
class Box422 {
    int value;
    int weight;
    Box422? next;

    function init(int v) {
        self.value = v + K22;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box422 n) {
        self.next = n;
    }
}
class Box423 {
    int value;
    int weight;
    Box423? next;

    function init(int v) {
        self.value = v + K23;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box423 n) {
        self.next = n;
    }
}
class Box424 {
    int value;
    int weight;
    Box424? next;

    function init(int v) {
        self.value = v + K24;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box424 n) {
        self.next = n;
    }
}
class Box425 {
    int value;
    int weight;
    Box425? next;

    function init(int v) {
        self.value = v + K25;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box425 n) {
        self.next = n;
    }
}
class Box426 {
    int value;
    int weight;
    Box426? next;

    function init(int v) {
        self.value = v + K26;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box426 n) {
        self.next = n;
    }
}
class Box427 {
    int value;
    int weight;
    Box427? next;

    function init(int v) {
        self.value = v + K27;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box427 n) {
        self.next = n;
    }
}
class Box428 {
    int value;
    int weight;
    Box428? next;

    function init(int v) {
        self.value = v + K28;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box428 n) {
        self.next = n;
    }
}
class Box429 {
    int value;
    int weight;
    Box429? next;

    function init(int v) {
        self.value = v + K29;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box429 n) {
        self.next = n;
    }
}
class Box430 {
    int value;
    int weight;
    Box430? next;

    function init(int v) {
        self.value = v + K30;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box430 n) {
        self.next = n;
    }
}
class Box431 {
    int value;
    int weight;
    Box431? next;

    function init(int v) {
        self.value = v + K31;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box431 n) {
        self.next = n;
    }
}
class Box432 {
    int value;
    int weight;
    Box432? next;

    function init(int v) {
        self.value = v + K32;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box432 n) {
        self.next = n;
    }
}
class Box433 {
    int value;
    int weight;
    Box433? next;

    function init(int v) {
        self.value = v + K33;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box433 n) {
        self.next = n;
    }
}
class Box434 {
    int value;
    int weight;
    Box434? next;

    function init(int v) {
        self.value = v + K34;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box434 n) {
        self.next = n;
    }
}
class Box435 {
    int value;
    int weight;
    Box435? next;

    function init(int v) {
        self.value = v + K35;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box435 n) {
        self.next = n;
    }
}
class Box436 {
    int value;
    int weight;
    Box436? next;

    function init(int v) {
        self.value = v + K36;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box436 n) {
        self.next = n;
    }
}
class Box437 {
    int value;
    int weight;
    Box437? next;

    function init(int v) {
        self.value = v + K37;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box437 n) {
        self.next = n;
    }
}
class Box438 {
    int value;
    int weight;
    Box438? next;

    function init(int v) {
        self.value = v + K38;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box438 n) {
        self.next = n;
    }
}
class Box439 {
    int value;
    int weight;
    Box439? next;

    function init(int v) {
        self.value = v + K39;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box439 n) {
        self.next = n;
    }
}
class Box440 {
    int value;
    int weight;
    Box440? next;

    function init(int v) {
        self.value = v + K40;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box440 n) {
        self.next = n;
    }
}
class Box441 {
    int value;
    int weight;
    Box441? next;

    function init(int v) {
        self.value = v + K41;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box441 n) {
        self.next = n;
    }
}
class Box442 {
    int value;
    int weight;
    Box442? next;

    function init(int v) {
        self.value = v + K42;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box442 n) {
        self.next = n;
    }
}
class Box443 {
    int value;
    int weight;
    Box443? next;

    function init(int v) {
        self.value = v + K43;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box443 n) {
        self.next = n;
    }
}
class Box444 {
    int value;
    int weight;
    Box444? next;

    function init(int v) {
        self.value = v + K44;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box444 n) {
        self.next = n;
    }
}
class Box445 {
    int value;
    int weight;
    Box445? next;

    function init(int v) {
        self.value = v + K45;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box445 n) {
        self.next = n;
    }
}
class Box446 {
    int value;
    int weight;
    Box446? next;

    function init(int v) {
        self.value = v + K46;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box446 n) {
        self.next = n;
    }
}
class Box447 {
    int value;
    int weight;
    Box447? next;

    function init(int v) {
        self.value = v + K47;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box447 n) {
        self.next = n;
    }
}
class Box448 {
    int value;
    int weight;
    Box448? next;

    function init(int v) {
        self.value = v + K48;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box448 n) {
        self.next = n;
    }
}
class Box449 {
    int value;
    int weight;
    Box449? next;

    function init(int v) {
        self.value = v + K49;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box449 n) {
        self.next = n;
    }
}
class Box450 {
    int value;
    int weight;
    Box450? next;

    function init(int v) {
        self.value = v + K50;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box450 n) {
        self.next = n;
    }
}
class Box451 {
    int value;
    int weight;
    Box451? next;

    function init(int v) {
        self.value = v + K51;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box451 n) {
        self.next = n;
    }
}
class Box452 {
    int value;
    int weight;
    Box452? next;

    function init(int v) {
        self.value = v + K52;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box452 n) {
        self.next = n;
    }
}
class Box453 {
    int value;
    int weight;
    Box453? next;

    function init(int v) {
        self.value = v + K53;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box453 n) {
        self.next = n;
    }
}
class Box454 {
    int value;
    int weight;
    Box454? next;

    function init(int v) {
        self.value = v + K54;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box454 n) {
        self.next = n;
    }
}
class Box455 {
    int value;
    int weight;
    Box455? next;

    function init(int v) {
        self.value = v + K55;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box455 n) {
        self.next = n;
    }
}
class Box456 {
    int value;
    int weight;
    Box456? next;

    function init(int v) {
        self.value = v + K56;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box456 n) {
        self.next = n;
    }
}
class Box457 {
    int value;
    int weight;
    Box457? next;

    function init(int v) {
        self.value = v + K57;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box457 n) {
        self.next = n;
    }
}
class Box458 {
    int value;
    int weight;
    Box458? next;

    function init(int v) {
        self.value = v + K58;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box458 n) {
        self.next = n;
    }
}
class Box459 {
    int value;
    int weight;
    Box459? next;

    function init(int v) {
        self.value = v + K59;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box459 n) {
        self.next = n;
    }
}
class Box460 {
    int value;
    int weight;
    Box460? next;

    function init(int v) {
        self.value = v + K60;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box460 n) {
        self.next = n;
    }
}
class Box461 {
    int value;
    int weight;
    Box461? next;

    function init(int v) {
        self.value = v + K61;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box461 n) {
        self.next = n;
    }
}
class Box462 {
    int value;
    int weight;
    Box462? next;

    function init(int v) {
        self.value = v + K62;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box462 n) {
        self.next = n;
    }
}
class Box463 {
    int value;
    int weight;
    Box463? next;

    function init(int v) {
        self.value = v + K63;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box463 n) {
        self.next = n;
    }
}
class Box464 {
    int value;
    int weight;
    Box464? next;

    function init(int v) {
        self.value = v + K64;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box464 n) {
        self.next = n;
    }
}
class Box465 {
    int value;
    int weight;
    Box465? next;

    function init(int v) {
        self.value = v + K65;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box465 n) {
        self.next = n;
    }
}
class Box466 {
    int value;
    int weight;
    Box466? next;

    function init(int v) {
        self.value = v + K66;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box466 n) {
        self.next = n;
    }
}
class Box467 {
    int value;
    int weight;
    Box467? next;

    function init(int v) {
        self.value = v + K67;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box467 n) {
        self.next = n;
    }
}
class Box468 {
    int value;
    int weight;
    Box468? next;

    function init(int v) {
        self.value = v + K68;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box468 n) {
        self.next = n;
    }
}
class Box469 {
    int value;
    int weight;
    Box469? next;

    function init(int v) {
        self.value = v + K69;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box469 n) {
        self.next = n;
    }
}
class Box470 {
    int value;
    int weight;
    Box470? next;

    function init(int v) {
        self.value = v + K70;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box470 n) {
        self.next = n;
    }
}
class Box471 {
    int value;
    int weight;
    Box471? next;

    function init(int v) {
        self.value = v + K71;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box471 n) {
        self.next = n;
    }
}
class Box472 {
    int value;
    int weight;
    Box472? next;

    function init(int v) {
        self.value = v + K72;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box472 n) {
        self.next = n;
    }
}
class Box473 {
    int value;
    int weight;
    Box473? next;

    function init(int v) {
        self.value = v + K73;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box473 n) {
        self.next = n;
    }
}
class Box474 {
    int value;
    int weight;
    Box474? next;

    function init(int v) {
        self.value = v + K74;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box474 n) {
        self.next = n;
    }
}
class Box475 {
    int value;
    int weight;
    Box475? next;

    function init(int v) {
        self.value = v + K75;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box475 n) {
        self.next = n;
    }
}
class Box476 {
    int value;
    int weight;
    Box476? next;

    function init(int v) {
        self.value = v + K76;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box476 n) {
        self.next = n;
    }
}
class Box477 {
    int value;
    int weight;
    Box477? next;

    function init(int v) {
        self.value = v + K77;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box477 n) {
        self.next = n;
    }
}
class Box478 {
    int value;
    int weight;
    Box478? next;

    function init(int v) {
        self.value = v + K78;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box478 n) {
        self.next = n;
    }
}
class Box479 {
    int value;
    int weight;
    Box479? next;

    function init(int v) {
        self.value = v + K79;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box479 n) {
        self.next = n;
    }
}
class Box480 {
    int value;
    int weight;
    Box480? next;

    function init(int v) {
        self.value = v + K80;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box480 n) {
        self.next = n;
    }
}
class Box481 {
    int value;
    int weight;
    Box481? next;

    function init(int v) {
        self.value = v + K81;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box481 n) {
        self.next = n;
    }
}
class Box482 {
    int value;
    int weight;
    Box482? next;

    function init(int v) {
        self.value = v + K82;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box482 n) {
        self.next = n;
    }
}
class Box483 {
    int value;
    int weight;
    Box483? next;

    function init(int v) {
        self.value = v + K83;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box483 n) {
        self.next = n;
    }
}
class Box484 {
    int value;
    int weight;
    Box484? next;

    function init(int v) {
        self.value = v + K84;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box484 n) {
        self.next = n;
    }
}
class Box485 {
    int value;
    int weight;
    Box485? next;

    function init(int v) {
        self.value = v + K85;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box485 n) {
        self.next = n;
    }
}
class Box486 {
    int value;
    int weight;
    Box486? next;

    function init(int v) {
        self.value = v + K86;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box486 n) {
        self.next = n;
    }
}
class Box487 {
    int value;
    int weight;
    Box487? next;

    function init(int v) {
        self.value = v + K87;
        self.weight = (v * 7) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box487 n) {
        self.next = n;
    }
}
class Box488 {
    int value;
    int weight;
    Box488? next;

    function init(int v) {
        self.value = v + K88;
        self.weight = (v * 8) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box488 n) {
        self.next = n;
    }
}
class Box489 {
    int value;
    int weight;
    Box489? next;

    function init(int v) {
        self.value = v + K89;
        self.weight = (v * 9) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box489 n) {
        self.next = n;
    }
}
class Box490 {
    int value;
    int weight;
    Box490? next;

    function init(int v) {
        self.value = v + K90;
        self.weight = (v * 10) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box490 n) {
        self.next = n;
    }
}
class Box491 {
    int value;
    int weight;
    Box491? next;

    function init(int v) {
        self.value = v + K91;
        self.weight = (v * 11) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box491 n) {
        self.next = n;
    }
}
class Box492 {
    int value;
    int weight;
    Box492? next;

    function init(int v) {
        self.value = v + K92;
        self.weight = (v * 12) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box492 n) {
        self.next = n;
    }
}
class Box493 {
    int value;
    int weight;
    Box493? next;

    function init(int v) {
        self.value = v + K93;
        self.weight = (v * 13) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 4) % 100000;
        }
        return self.weight;
    }

    function setNext(Box493 n) {
        self.next = n;
    }
}
class Box494 {
    int value;
    int weight;
    Box494? next;

    function init(int v) {
        self.value = v + K94;
        self.weight = (v * 1) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 5) % 100000;
        }
        return self.weight;
    }

    function setNext(Box494 n) {
        self.next = n;
    }
}
class Box495 {
    int value;
    int weight;
    Box495? next;

    function init(int v) {
        self.value = v + K95;
        self.weight = (v * 2) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 6) % 100000;
        }
        return self.weight;
    }

    function setNext(Box495 n) {
        self.next = n;
    }
}
class Box496 {
    int value;
    int weight;
    Box496? next;

    function init(int v) {
        self.value = v + K96;
        self.weight = (v * 3) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 7) % 100000;
        }
        return self.weight;
    }

    function setNext(Box496 n) {
        self.next = n;
    }
}
class Box497 {
    int value;
    int weight;
    Box497? next;

    function init(int v) {
        self.value = v + K97;
        self.weight = (v * 4) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 1) % 100000;
        }
        return self.weight;
    }

    function setNext(Box497 n) {
        self.next = n;
    }
}
class Box498 {
    int value;
    int weight;
    Box498? next;

    function init(int v) {
        self.value = v + K98;
        self.weight = (v * 5) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 2) % 100000;
        }
        return self.weight;
    }

    function setNext(Box498 n) {
        self.next = n;
    }
}
class Box499 {
    int value;
    int weight;
    Box499? next;

    function init(int v) {
        self.value = v + K99;
        self.weight = (v * 6) % 100000;
        self.next = ();
    }

    function compute(int salt) returns int {
        int? lifted = self.value + salt;
        if lifted is int {
            return (lifted * 3) % 100000;
        }
        return self.weight;
    }

    function setNext(Box499 n) {
        self.next = n;
    }
}

public function main() {
    int total = 0;
    Box0 b0 = new Box0(1);
    total = (total + b0.compute(2)) % 1000000;
    Box1 b1 = new Box1(2);
    total = (total + b1.compute(3)) % 1000000;
    Box2 b2 = new Box2(3);
    total = (total + b2.compute(4)) % 1000000;
    Box3 b3 = new Box3(4);
    total = (total + b3.compute(5)) % 1000000;
    Box4 b4 = new Box4(5);
    total = (total + b4.compute(6)) % 1000000;
    Box5 b5 = new Box5(6);
    total = (total + b5.compute(7)) % 1000000;
    Box6 b6 = new Box6(7);
    total = (total + b6.compute(8)) % 1000000;
    Box7 b7 = new Box7(8);
    total = (total + b7.compute(9)) % 1000000;
    Box8 b8 = new Box8(9);
    total = (total + b8.compute(10)) % 1000000;
    Box9 b9 = new Box9(10);
    total = (total + b9.compute(11)) % 1000000;
    Box10 b10 = new Box10(11);
    total = (total + b10.compute(12)) % 1000000;
    Box11 b11 = new Box11(12);
    total = (total + b11.compute(13)) % 1000000;
    Box12 b12 = new Box12(13);
    total = (total + b12.compute(14)) % 1000000;
    Box13 b13 = new Box13(14);
    total = (total + b13.compute(15)) % 1000000;
    Box14 b14 = new Box14(15);
    total = (total + b14.compute(16)) % 1000000;
    Box15 b15 = new Box15(16);
    total = (total + b15.compute(17)) % 1000000;
    Box16 b16 = new Box16(17);
    total = (total + b16.compute(18)) % 1000000;
    Box17 b17 = new Box17(18);
    total = (total + b17.compute(19)) % 1000000;
    Box18 b18 = new Box18(19);
    total = (total + b18.compute(20)) % 1000000;
    Box19 b19 = new Box19(20);
    total = (total + b19.compute(21)) % 1000000;
    Box20 b20 = new Box20(21);
    total = (total + b20.compute(22)) % 1000000;
    Box21 b21 = new Box21(22);
    total = (total + b21.compute(23)) % 1000000;
    Box22 b22 = new Box22(23);
    total = (total + b22.compute(24)) % 1000000;
    Box23 b23 = new Box23(24);
    total = (total + b23.compute(25)) % 1000000;
    Box24 b24 = new Box24(25);
    total = (total + b24.compute(26)) % 1000000;
    Box25 b25 = new Box25(26);
    total = (total + b25.compute(27)) % 1000000;
    Box26 b26 = new Box26(27);
    total = (total + b26.compute(28)) % 1000000;
    Box27 b27 = new Box27(28);
    total = (total + b27.compute(29)) % 1000000;
    Box28 b28 = new Box28(29);
    total = (total + b28.compute(30)) % 1000000;
    Box29 b29 = new Box29(30);
    total = (total + b29.compute(31)) % 1000000;
    Box30 b30 = new Box30(31);
    total = (total + b30.compute(32)) % 1000000;
    Box31 b31 = new Box31(32);
    total = (total + b31.compute(33)) % 1000000;
    Box32 b32 = new Box32(33);
    total = (total + b32.compute(34)) % 1000000;
    Box33 b33 = new Box33(34);
    total = (total + b33.compute(35)) % 1000000;
    Box34 b34 = new Box34(35);
    total = (total + b34.compute(36)) % 1000000;
    Box35 b35 = new Box35(36);
    total = (total + b35.compute(37)) % 1000000;
    Box36 b36 = new Box36(37);
    total = (total + b36.compute(38)) % 1000000;
    Box37 b37 = new Box37(38);
    total = (total + b37.compute(39)) % 1000000;
    Box38 b38 = new Box38(39);
    total = (total + b38.compute(40)) % 1000000;
    Box39 b39 = new Box39(40);
    total = (total + b39.compute(41)) % 1000000;
    Box40 b40 = new Box40(41);
    total = (total + b40.compute(42)) % 1000000;
    Box41 b41 = new Box41(42);
    total = (total + b41.compute(43)) % 1000000;
    Box42 b42 = new Box42(43);
    total = (total + b42.compute(44)) % 1000000;
    Box43 b43 = new Box43(44);
    total = (total + b43.compute(45)) % 1000000;
    Box44 b44 = new Box44(45);
    total = (total + b44.compute(46)) % 1000000;
    Box45 b45 = new Box45(46);
    total = (total + b45.compute(47)) % 1000000;
    Box46 b46 = new Box46(47);
    total = (total + b46.compute(48)) % 1000000;
    Box47 b47 = new Box47(48);
    total = (total + b47.compute(49)) % 1000000;
    Box48 b48 = new Box48(49);
    total = (total + b48.compute(50)) % 1000000;
    Box49 b49 = new Box49(50);
    total = (total + b49.compute(51)) % 1000000;
    Box50 b50 = new Box50(51);
    total = (total + b50.compute(52)) % 1000000;
    Box51 b51 = new Box51(52);
    total = (total + b51.compute(53)) % 1000000;
    Box52 b52 = new Box52(53);
    total = (total + b52.compute(54)) % 1000000;
    Box53 b53 = new Box53(54);
    total = (total + b53.compute(55)) % 1000000;
    Box54 b54 = new Box54(55);
    total = (total + b54.compute(56)) % 1000000;
    Box55 b55 = new Box55(56);
    total = (total + b55.compute(57)) % 1000000;
    Box56 b56 = new Box56(57);
    total = (total + b56.compute(58)) % 1000000;
    Box57 b57 = new Box57(58);
    total = (total + b57.compute(59)) % 1000000;
    Box58 b58 = new Box58(59);
    total = (total + b58.compute(60)) % 1000000;
    Box59 b59 = new Box59(60);
    total = (total + b59.compute(61)) % 1000000;
    Box60 b60 = new Box60(61);
    total = (total + b60.compute(62)) % 1000000;
    Box61 b61 = new Box61(62);
    total = (total + b61.compute(63)) % 1000000;
    Box62 b62 = new Box62(63);
    total = (total + b62.compute(64)) % 1000000;
    Box63 b63 = new Box63(64);
    total = (total + b63.compute(65)) % 1000000;
    Box64 b64 = new Box64(65);
    total = (total + b64.compute(66)) % 1000000;
    Box65 b65 = new Box65(66);
    total = (total + b65.compute(67)) % 1000000;
    Box66 b66 = new Box66(67);
    total = (total + b66.compute(68)) % 1000000;
    Box67 b67 = new Box67(68);
    total = (total + b67.compute(69)) % 1000000;
    Box68 b68 = new Box68(69);
    total = (total + b68.compute(70)) % 1000000;
    Box69 b69 = new Box69(70);
    total = (total + b69.compute(71)) % 1000000;
    Box70 b70 = new Box70(71);
    total = (total + b70.compute(72)) % 1000000;
    Box71 b71 = new Box71(72);
    total = (total + b71.compute(73)) % 1000000;
    Box72 b72 = new Box72(73);
    total = (total + b72.compute(74)) % 1000000;
    Box73 b73 = new Box73(74);
    total = (total + b73.compute(75)) % 1000000;
    Box74 b74 = new Box74(75);
    total = (total + b74.compute(76)) % 1000000;
    Box75 b75 = new Box75(76);
    total = (total + b75.compute(77)) % 1000000;
    Box76 b76 = new Box76(77);
    total = (total + b76.compute(78)) % 1000000;
    Box77 b77 = new Box77(78);
    total = (total + b77.compute(79)) % 1000000;
    Box78 b78 = new Box78(79);
    total = (total + b78.compute(80)) % 1000000;
    Box79 b79 = new Box79(80);
    total = (total + b79.compute(81)) % 1000000;
    Box0 c0 = new Box0(11);
    Box0 c1 = new Box0(22);
    c0.setNext(c1);
    total = (total + c0.compute(3) + c1.compute(5)) % 1000000;
    io:println(total);
}
