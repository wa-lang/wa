// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package loong64

// 寄存器名字列表
var _Register = []string{
	// 通用寄存器
	REG_R0:  "R0",
	REG_R1:  "R1",
	REG_R2:  "R2",
	REG_R3:  "R3",
	REG_R4:  "R4",
	REG_R5:  "R5",
	REG_R6:  "R6",
	REG_R7:  "R7",
	REG_R8:  "R8",
	REG_R9:  "R9",
	REG_R10: "R10",
	REG_R11: "R11",
	REG_R12: "R12",
	REG_R13: "R13",
	REG_R14: "R14",
	REG_R15: "R15",
	REG_R16: "R16",
	REG_R17: "R17",
	REG_R18: "R18",
	REG_R19: "R19",
	REG_R20: "R20",
	REG_R21: "R21",
	REG_R22: "R22",
	REG_R23: "R23",
	REG_R24: "R24",
	REG_R25: "R25",
	REG_R26: "R26",
	REG_R27: "R27",
	REG_R28: "R28",
	REG_R29: "R29",
	REG_R30: "R30",
	REG_R31: "R31",

	// 浮点数寄存器(F/D扩展)
	REG_F0:  "F0",
	REG_F1:  "F1",
	REG_F2:  "F2",
	REG_F3:  "F3",
	REG_F4:  "F4",
	REG_F5:  "F5",
	REG_F6:  "F6",
	REG_F7:  "F7",
	REG_F8:  "F8",
	REG_F9:  "F9",
	REG_F10: "F10",
	REG_F11: "F11",
	REG_F12: "F12",
	REG_F13: "F13",
	REG_F14: "F14",
	REG_F15: "F15",
	REG_F16: "F16",
	REG_F17: "F17",
	REG_F18: "F18",
	REG_F19: "F19",
	REG_F20: "F20",
	REG_F21: "F21",
	REG_F22: "F22",
	REG_F23: "F23",
	REG_F24: "F24",
	REG_F25: "F25",
	REG_F26: "F26",
	REG_F27: "F27",
	REG_F28: "F28",
	REG_F29: "F29",
	REG_F30: "F30",
	REG_F31: "F31",

	// LSX
	REG_VF0:  "VF0",
	REG_VF1:  "VF1",
	REG_VF2:  "VF2",
	REG_VF3:  "VF3",
	REG_VF4:  "VF4",
	REG_VF5:  "VF5",
	REG_VF6:  "VF6",
	REG_VF7:  "VF7",
	REG_VF8:  "VF8",
	REG_VF9:  "VF9",
	REG_VF10: "VF10",
	REG_VF11: "VF11",
	REG_VF12: "VF12",
	REG_VF13: "VF13",
	REG_VF14: "VF14",
	REG_VF15: "VF15",
	REG_VF16: "VF16",
	REG_VF17: "VF17",
	REG_VF18: "VF18",
	REG_VF19: "VF19",
	REG_VF20: "VF20",
	REG_VF21: "VF21",
	REG_VF22: "VF22",
	REG_VF23: "VF23",
	REG_VF24: "VF24",
	REG_VF25: "VF25",
	REG_VF26: "VF26",
	REG_VF27: "VF27",
	REG_VF28: "VF28",
	REG_VF29: "VF29",
	REG_VF30: "VF30",
	REG_VF31: "VF31",

	// LASX
	REG_XF0:  "XF0",
	REG_XF1:  "XF1",
	REG_XF2:  "XF2",
	REG_XF3:  "XF3",
	REG_XF4:  "XF4",
	REG_XF5:  "XF5",
	REG_XF6:  "XF6",
	REG_XF7:  "XF7",
	REG_XF8:  "XF8",
	REG_XF9:  "XF9",
	REG_XF10: "XF10",
	REG_XF11: "XF11",
	REG_XF12: "XF12",
	REG_XF13: "XF13",
	REG_XF14: "XF14",
	REG_XF15: "XF15",
	REG_XF16: "XF16",
	REG_XF17: "XF17",
	REG_XF18: "XF18",
	REG_XF19: "XF19",
	REG_XF20: "XF20",
	REG_XF21: "XF21",
	REG_XF22: "XF22",
	REG_XF23: "XF23",
	REG_XF24: "XF24",
	REG_XF25: "XF25",
	REG_XF26: "XF26",
	REG_XF27: "XF27",
	REG_XF28: "XF28",
	REG_XF29: "XF29",
	REG_XF30: "XF30",
	REG_XF31: "XF31",
}

// 寄存器名字列表(中文)
var _ZhRegister = []string{
	// 通用寄存器
	REG_R0:  "甲格",
	REG_R1:  "乙格",
	REG_R2:  "丙格",
	REG_R3:  "丁格",
	REG_R4:  "戊格",
	REG_R5:  "己格",
	REG_R6:  "庚格",
	REG_R7:  "辛格",
	REG_R8:  "壬格",
	REG_R9:  "癸格",
	REG_R10: "子格",
	REG_R11: "丑格",
	REG_R12: "寅格",
	REG_R13: "卯格",
	REG_R14: "辰格",
	REG_R15: "巳格",
	REG_R16: "午格",
	REG_R17: "未格",
	REG_R18: "申格",
	REG_R19: "酉格",
	REG_R20: "戌格",
	REG_R21: "亥格",
	REG_R22: "乾格",
	REG_R23: "坤格",
	REG_R24: "震格",
	REG_R25: "巽格",
	REG_R26: "坎格",
	REG_R27: "离格",
	REG_R28: "艮格",
	REG_R29: "兑格",
	REG_R30: "阴格",
	REG_R31: "阳格",

	// 浮点数寄存器(F/D扩展)
	REG_F0:  "甲皿",
	REG_F1:  "乙皿",
	REG_F2:  "丙皿",
	REG_F3:  "丁皿",
	REG_F4:  "戊皿",
	REG_F5:  "己皿",
	REG_F6:  "庚皿",
	REG_F7:  "辛皿",
	REG_F8:  "壬皿",
	REG_F9:  "癸皿",
	REG_F10: "子皿",
	REG_F11: "丑皿",
	REG_F12: "寅皿",
	REG_F13: "卯皿",
	REG_F14: "辰皿",
	REG_F15: "巳皿",
	REG_F16: "午皿",
	REG_F17: "未皿",
	REG_F18: "申皿",
	REG_F19: "酉皿",
	REG_F20: "戌皿",
	REG_F21: "亥皿",
	REG_F22: "乾皿",
	REG_F23: "坤皿",
	REG_F24: "震皿",
	REG_F25: "巽皿",
	REG_F26: "坎皿",
	REG_F27: "离皿",
	REG_F28: "艮皿",
	REG_F29: "兑皿",
	REG_F30: "阴皿",
	REG_F31: "阳皿",

	// LSX
	REG_VF0:  "威甲皿",
	REG_VF1:  "威乙皿",
	REG_VF2:  "威丙皿",
	REG_VF3:  "威丁皿",
	REG_VF4:  "威戊皿",
	REG_VF5:  "威己皿",
	REG_VF6:  "威庚皿",
	REG_VF7:  "威辛皿",
	REG_VF8:  "威壬皿",
	REG_VF9:  "威癸皿",
	REG_VF10: "威子皿",
	REG_VF11: "威丑皿",
	REG_VF12: "威寅皿",
	REG_VF13: "威卯皿",
	REG_VF14: "威辰皿",
	REG_VF15: "威巳皿",
	REG_VF16: "威午皿",
	REG_VF17: "威未皿",
	REG_VF18: "威申皿",
	REG_VF19: "威酉皿",
	REG_VF20: "威戌皿",
	REG_VF21: "威亥皿",
	REG_VF22: "威乾皿",
	REG_VF23: "威坤皿",
	REG_VF24: "威震皿",
	REG_VF25: "威巽皿",
	REG_VF26: "威坎皿",
	REG_VF27: "威离皿",
	REG_VF28: "威艮皿",
	REG_VF29: "威兑皿",
	REG_VF30: "威阴皿",
	REG_VF31: "威阳皿",

	// LASX
	REG_XF0:  "叉甲皿",
	REG_XF1:  "叉乙皿",
	REG_XF2:  "叉丙皿",
	REG_XF3:  "叉丁皿",
	REG_XF4:  "叉戊皿",
	REG_XF5:  "叉己皿",
	REG_XF6:  "叉庚皿",
	REG_XF7:  "叉辛皿",
	REG_XF8:  "叉壬皿",
	REG_XF9:  "叉癸皿",
	REG_XF10: "叉子皿",
	REG_XF11: "叉丑皿",
	REG_XF12: "叉寅皿",
	REG_XF13: "叉卯皿",
	REG_XF14: "叉辰皿",
	REG_XF15: "叉巳皿",
	REG_XF16: "叉午皿",
	REG_XF17: "叉未皿",
	REG_XF18: "叉申皿",
	REG_XF19: "叉酉皿",
	REG_XF20: "叉戌皿",
	REG_XF21: "叉亥皿",
	REG_XF22: "叉乾皿",
	REG_XF23: "叉坤皿",
	REG_XF24: "叉震皿",
	REG_XF25: "叉巽皿",
	REG_XF26: "叉坎皿",
	REG_XF27: "叉离皿",
	REG_XF28: "叉艮皿",
	REG_XF29: "叉兑皿",
	REG_XF30: "叉阴皿",
	REG_XF31: "叉阳皿",
}

// 寄存器别名
var _RegisterAlias = []string{
	REG_ZERO: "ZERO", // 零寄存器
	REG_RA:   "RA",   // 返回地址
	REG_TP:   "TP",   // 线程指针
	REG_SP:   "SP",   // 栈指针
	REG_A0:   "A0",   // 函数参数/返回值
	REG_A1:   "A1",   // 函数参数/返回值
	REG_A2:   "A2",   // 函数参数
	REG_A3:   "A3",   // 函数参数
	REG_A4:   "A4",   // 函数参数
	REG_A5:   "A5",   // 函数参数
	REG_A6:   "A6",   // 函数参数
	REG_A7:   "A7",   // 函数参数
	REG_T0:   "T0",   // 临时变量
	REG_T1:   "T1",   // 临时变量
	REG_T2:   "T2",   // 临时变量
	REG_T3:   "T3",   // 临时变量
	REG_T4:   "T4",   // 临时变量
	REG_T5:   "T5",   // 临时变量
	REG_T6:   "T6",   // 临时变量
	REG_T7:   "T7",   // 临时变量
	REG_T8:   "T8",   // 临时变量
	REG_X:    "X",    // 保留寄存器
	REG_FP:   "FP",   // 帧指针
	REG_S0:   "S0",   // Saved register
	REG_S1:   "S1",   // Saved register
	REG_S2:   "S2",   // Saved register
	REG_S3:   "S3",   // Saved register
	REG_S4:   "S4",   // Saved register
	REG_S5:   "S5",   // Saved register
	REG_S6:   "S6",   // Saved register
	REG_S7:   "S7",   // Saved register
	REG_S8:   "S8",   // Saved register

	REG_FA0:  "FA0",  // 函数参数/返回值
	REG_FA1:  "FA1",  // 函数参数/返回值
	REG_FA2:  "FA2",  // 函数参数
	REG_FA3:  "FA3",  // 函数参数
	REG_FA4:  "FA4",  // 函数参数
	REG_FA5:  "FA5",  // 函数参数
	REG_FA6:  "FA6",  // 函数参数
	REG_FA7:  "FA7",  // 函数参数
	REG_FT0:  "FT0",  // 临时变量
	REG_FT1:  "FT1",  // 临时变量
	REG_FT2:  "FT2",  // 临时变量
	REG_FT3:  "FT3",  // 临时变量
	REG_FT4:  "FT4",  // 临时变量
	REG_FT5:  "FT5",  // 临时变量
	REG_FT6:  "FT6",  // 临时变量
	REG_FT7:  "FT7",  // 临时变量
	REG_FT8:  "FT8",  // 临时变量
	REG_FT9:  "FT9",  // 临时变量
	REG_FT10: "FT10", // 临时变量
	REG_FT11: "FT11", // 临时变量
	REG_FT12: "FT12", // 临时变量
	REG_FT13: "FT13", // 临时变量
	REG_FT14: "FT14", // 临时变量
	REG_FT15: "FT15", // 临时变量
	REG_FS0:  "FS0",  // Saved register
	REG_FS1:  "FS1",  // Saved register
	REG_FS2:  "FS2",  // Saved register
	REG_FS3:  "FS3",  // Saved register
	REG_FS4:  "FS4",  // Saved register
	REG_FS5:  "FS5",  // Saved register
	REG_FS6:  "FS6",  // Saved register
	REG_FS7:  "FS7",  // Saved register
}

var _ZhRegisterAlias = []string{
	REG_ZERO: "零格",  // 零寄存器
	REG_RA:   "回格",  // 返回地址
	REG_TP:   "线格",  // 线程指针
	REG_SP:   "栈格",  // 栈指针
	REG_FP:   "帧格",  // 帧指针
	REG_A0:   "参甲格", // 函数参数/返回值
	REG_A1:   "参乙格", // 函数参数/返回值
	REG_A2:   "参丙格", // 函数参数
	REG_A3:   "参丁格", // 函数参数
	REG_A4:   "参戊格", // 函数参数
	REG_A5:   "参己格", // 函数参数
	REG_A6:   "参庚格", // 函数参数
	REG_A7:   "参辛格", // 函数参数
	REG_T0:   "暂甲格", // 临时变量
	REG_T1:   "暂乙格", // 临时变量
	REG_T2:   "暂丙格", // 临时变量
	REG_T3:   "暂丁格", // 临时变量
	REG_T4:   "暂戊格", // 临时变量
	REG_T5:   "暂己格", // 临时变量
	REG_T6:   "暂庚格", // 临时变量
	REG_T7:   "暂辛格", // 临时变量
	REG_T8:   "暂壬格", // 临时变量
	REG_X:    "保留格", // 保留寄存器
	REG_S0:   "守甲格", // Saved register
	REG_S1:   "守乙格", // Saved register
	REG_S2:   "守丙格", // Saved register
	REG_S3:   "守丁格", // Saved register
	REG_S4:   "守戊格", // Saved register
	REG_S5:   "守己格", // Saved register
	REG_S6:   "守庚格", // Saved register
	REG_S7:   "守辛格", // Saved register
	REG_S8:   "守壬格", // Saved register

	REG_FA0:  "参甲皿", // 函数参数/返回值
	REG_FA1:  "参乙皿", // 函数参数/返回值
	REG_FA2:  "参丙皿", // 函数参数
	REG_FA3:  "参丁皿", // 函数参数
	REG_FA4:  "参戊皿", // 函数参数
	REG_FA5:  "参己皿", // 函数参数
	REG_FA6:  "参庚皿", // 函数参数
	REG_FA7:  "参辛皿", // 函数参数
	REG_FT0:  "暂甲皿", // 临时变量
	REG_FT1:  "暂乙皿", // 临时变量
	REG_FT2:  "暂丙皿", // 临时变量
	REG_FT3:  "暂丁皿", // 临时变量
	REG_FT4:  "暂戊皿", // 临时变量
	REG_FT5:  "暂己皿", // 临时变量
	REG_FT6:  "暂庚皿", // 临时变量
	REG_FT7:  "暂辛皿", // 临时变量
	REG_FT8:  "暂壬皿", // 临时变量
	REG_FT9:  "暂癸皿", // 临时变量
	REG_FT10: "暂子皿", // 临时变量
	REG_FT11: "暂丑皿", // 临时变量
	REG_FT12: "暂寅皿", // 临时变量
	REG_FT13: "暂卯皿", // 临时变量
	REG_FT14: "暂辰皿", // 临时变量
	REG_FT15: "暂巳皿", // 临时变量
	REG_FS0:  "守甲皿", // Saved register
	REG_FS1:  "守乙皿", // Saved register
	REG_FS2:  "守丙皿", // Saved register
	REG_FS3:  "守丁皿", // Saved register
	REG_FS4:  "守戊皿", // Saved register
	REG_FS5:  "守己皿", // Saved register
	REG_FS6:  "守庚皿", // Saved register
	REG_FS7:  "守辛皿", // Saved register
}

// 中文指令集
//
// 整数的后缀: 微/短/字/长/正, 对应 B/H/W/D/U
// 浮点数的后缀: 半/单/双/字/长/正, 对应 H/S/D/W/L/U
// 立即数: 立 对应 I
// 程序计数器: 计 对应 PC
// 其他指令部分按字面意思翻译, 比如 加 对应 ADD, 减 对应 SUB
var _ZhAnames = []string{
	AADDI_D:       "加立.长",
	AADDI_W:       "加立.字",
	AADDU16I_D:    "加正16立.长",
	AADD_D:        "加.长",
	AADD_W:        "加.字",
	AALSL_D:       "左移加.长",
	AALSL_W:       "左移加.字",
	AALSL_WU:      "左移加.字正",
	AAMADD_B:      "原子加乘.微",
	AAMADD_D:      "原子加乘.长",
	AAMADD_DB_B:   "原子加乘.双缓冲.微",
	AAMADD_DB_D:   "原子加乘.双缓冲.长",
	AAMADD_DB_H:   "原子加乘.双缓冲.短",
	AAMADD_DB_W:   "原子加乘.双缓冲.字",
	AAMADD_H:      "原子加乘.短",
	AAMADD_W:      "原子加乘.字",
	AAMAND_D:      "原子与.长",
	AAMAND_DB_D:   "原子与.双缓冲.长",
	AAMAND_DB_W:   "原子与.双缓冲.字",
	AAMAND_W:      "原子与.字",
	AAMCAS_B:      "原子比较交换.微",
	AAMCAS_D:      "原子比较交换.长",
	AAMCAS_DB_B:   "原子比较交换.双缓冲.微",
	AAMCAS_DB_D:   "原子比较交换.双缓冲.长",
	AAMCAS_DB_H:   "原子比较交换.双缓冲.短",
	AAMCAS_DB_W:   "原子比较交换.双缓冲.字",
	AAMCAS_H:      "原子比较交换.短",
	AAMCAS_W:      "原子比较交换.字",
	AAMMAX_D:      "原子最大.长",
	AAMMAX_DB_D:   "原子最大.双缓冲.长",
	AAMMAX_DB_DU:  "原子最大.双缓冲.长正",
	AAMMAX_DB_W:   "原子最大.双缓冲.字",
	AAMMAX_DB_WU:  "原子最大.双缓冲.字正",
	AAMMAX_DU:     "原子最大.长正",
	AAMMAX_W:      "原子最大.字",
	AAMMAX_WU:     "原子最大.字正",
	AAMMIN_D:      "原子最小.长",
	AAMMIN_DB_D:   "原子最小.双缓冲.长",
	AAMMIN_DB_DU:  "原子最小.双缓冲.长正",
	AAMMIN_DB_W:   "原子最小.双缓冲.字",
	AAMMIN_DB_WU:  "原子最小.双缓冲.字正",
	AAMMIN_DU:     "原子最小.长正",
	AAMMIN_W:      "原子最小.字",
	AAMMIN_WU:     "原子最小.字正",
	AAMOR_D:       "原子或.长",
	AAMOR_DB_D:    "原子或.双缓冲.长",
	AAMOR_DB_W:    "原子或.双缓冲.字",
	AAMOR_W:       "原子或.字",
	AAMSWAP_B:     "原子交换.微",
	AAMSWAP_D:     "原子交换.长",
	AAMSWAP_DB_B:  "原子交换.双缓冲.微",
	AAMSWAP_DB_D:  "原子交换.双缓冲.长",
	AAMSWAP_DB_H:  "原子交换.双缓冲.短",
	AAMSWAP_DB_W:  "原子交换.双缓冲.字",
	AAMSWAP_H:     "原子交换.短",
	AAMSWAP_W:     "原子交换.字",
	AAMXOR_D:      "原子异或.长",
	AAMXOR_DB_D:   "原子异或.双缓冲.长",
	AAMXOR_DB_W:   "原子异或.双缓冲.字",
	AAMXOR_W:      "原子异或.字",
	AAND:          "与",
	AANDI:         "与立",
	AANDN:         "与非",
	AASRTGT_D:     "右移大于.长",
	AASRTLE_D:     "右移小于等于.长",
	AB:            "跳转",
	ABCEQZ:        "跳转.浮点数.等于零",
	ABCNEZ:        "跳转.浮点数.不等零",
	ABEQ:          "跳转.相等",
	ABEQZ:         "跳转.等于零",
	ABGE:          "跳转.大于等于",
	ABGEU:         "跳转.大于等于.正",
	ABITREV_4B:    "位反转.4微",
	ABITREV_8B:    "位反转.8微",
	ABITREV_D:     "位反转.长",
	ABITREV_W:     "位反转.字",
	ABL:           "跳转.可返",
	ABLT:          "跳转.可返.小于",
	ABLTU:         "跳转.可返.小于.正",
	ABNE:          "跳转.不相等",
	ABNEZ:         "跳转.不等零",
	ABREAK:        "断点",
	ABSTRINS_D:    "位域插入.长",
	ABSTRINS_W:    "位域插入.字",
	ABSTRPICK_D:   "位域提取.长",
	ABSTRPICK_W:   "位域提取.字",
	ABYTEPICK_D:   "字节提取.长",
	ABYTEPICK_W:   "字节提取.字",
	ACACOP:        "缓存操作",
	ACLO_D:        "前导1计数.长",
	ACLO_W:        "前导1计数.字",
	ACLZ_D:        "前导0计数.长",
	ACLZ_W:        "前导0计数.字",
	ACPUCFG:       "CPU配置",
	ACRCC_W_B_W:   "CRC校验累加.字.微.字",
	ACRCC_W_D_W:   "CRC校验累加.字.长.字",
	ACRCC_W_H_W:   "CRC校验累加.字.短.字",
	ACRCC_W_W_W:   "CRC校验累加.字.字.字",
	ACRC_W_B_W:    "CRC校验.字.微.字",
	ACRC_W_D_W:    "CRC校验.字.长.字",
	ACRC_W_H_W:    "CRC校验.字.短.字",
	ACRC_W_W_W:    "CRC校验.字.字.字",
	ACSRRD:        "读控制状态",
	ACSRWR:        "写控制状态",
	ACSRXCHG:      "交换控制状态",
	ACTO_D:        "后缀1计数.长",
	ACTO_W:        "后缀1计数.字",
	ACTZ_D:        "后缀0计数.长",
	ACTZ_W:        "后缀0计数.字",
	ADBAR:         "数据屏障",
	ADBCL:         "数据缓存清除",
	ADIV_D:        "除.长",
	ADIV_DU:       "除.长正",
	ADIV_W:        "除.字",
	ADIV_WU:       "除.字正",
	AERTN:         "异常返回",
	AEXT_W_B:      "符号扩展.字.微",
	AEXT_W_H:      "符号扩展.字.短",
	AFABS_D:       "浮绝对值.双",
	AFABS_S:       "浮绝对值.单",
	AFADD_D:       "浮加.双",
	AFADD_S:       "浮加.单",
	AFCLASS_D:     "浮分类.双",
	AFCLASS_S:     "浮分类.单",
	AFCMP_CAF_D:   "浮比较.条件位.假.双",
	AFCMP_CAF_S:   "浮比较.条件位.假.单",
	AFCMP_CEQ_D:   "浮比较.条件位.等.双",
	AFCMP_CEQ_S:   "浮比较.条件位.等.单",
	AFCMP_CLE_D:   "浮比较.条件位.小于等于.双",
	AFCMP_CLE_S:   "浮比较.条件位.小于等于.单",
	AFCMP_CLT_D:   "浮比较.条件位.小于.双",
	AFCMP_CLT_S:   "浮比较.条件位.小于.单",
	AFCMP_CNE_D:   "浮比较.条件位.不等.双",
	AFCMP_CNE_S:   "浮比较.条件位.不等.",
	AFCMP_COR_D:   "浮比较.条件位.真.双",
	AFCMP_COR_S:   "浮比较.条件位.真.单",
	AFCMP_CUEQ_D:  "浮比较.条件位.无序等.双",
	AFCMP_CUEQ_S:  "浮比较.条件位.无序等.单",
	AFCMP_CULE_D:  "浮比较.条件位.无序小于等于.双",
	AFCMP_CULE_S:  "浮比较.条件位.无序小于等于.单",
	AFCMP_CULT_D:  "浮比较.条件位.无序小于.双",
	AFCMP_CULT_S:  "浮比较.条件位.无序小于.单",
	AFCMP_CUNE_D:  "浮比较.条件位.无序不等.双",
	AFCMP_CUNE_S:  "浮比较.条件位.无序不等.单",
	AFCMP_CUN_D:   "浮比较.条件位.无序.双",
	AFCMP_CUN_S:   "浮比较.条件位.无序.单",
	AFCMP_SAF_D:   "浮比较.状态位.假.双",
	AFCMP_SAF_S:   "浮比较.状态位.假.单",
	AFCMP_SEQ_D:   "浮比较.状态位.等.双",
	AFCMP_SEQ_S:   "浮比较.状态位.等.单",
	AFCMP_SLE_D:   "浮比较.状态位.小于等于.双",
	AFCMP_SLE_S:   "浮比较.状态位.小于等于.单",
	AFCMP_SLT_D:   "浮比较.状态位.小于.双",
	AFCMP_SLT_S:   "浮比较.状态位.小于.单",
	AFCMP_SNE_D:   "浮比较.状态位.不等.双",
	AFCMP_SNE_S:   "浮比较.状态位.不等.单",
	AFCMP_SOR_D:   "浮比较.状态位.真.双",
	AFCMP_SOR_S:   "浮比较.状态位.真.单",
	AFCMP_SUEQ_D:  "浮比较.状态位.无序等.双",
	AFCMP_SUEQ_S:  "浮比较.状态位.无序等.单",
	AFCMP_SULE_D:  "浮比较.状态位.无序小于等于.双",
	AFCMP_SULE_S:  "浮比较.状态位.无序小于等于.单",
	AFCMP_SULT_D:  "浮比较.状态位.无序小于.双",
	AFCMP_SULT_S:  "浮比较.状态位.无序小于.单",
	AFCMP_SUNE_D:  "浮比较.状态位.无序不等.双",
	AFCMP_SUNE_S:  "浮比较.状态位.无序不等.单",
	AFCMP_SUN_D:   "浮比较.状态位.无序.双",
	AFCMP_SUN_S:   "浮比较.状态位.无序.单",
	AFCOPYSIGN_D:  "浮拷贝符号位.双",
	AFCOPYSIGN_S:  "浮拷贝符号位.单",
	AFCVT_D_S:     "浮转换.双.单",
	AFCVT_S_D:     "浮转换.单.双",
	AFDIV_D:       "浮除.双",
	AFDIV_S:       "浮除.单",
	AFFINT_D_L:    "浮转浮点整数.双.长",
	AFFINT_D_W:    "浮转浮点整数.双.字",
	AFFINT_S_L:    "浮转浮点整数.单.长",
	AFFINT_S_W:    "浮转浮点整数.单.字",
	AFLDGT_D:      "浮装载.大于.双",
	AFLDGT_S:      "浮装载.大于.单",
	AFLDLE_D:      "浮装载.小于等于.双",
	AFLDLE_S:      "浮装载.小于等于.单",
	AFLDX_D:       "浮装载变址.双",
	AFLDX_S:       "浮装载变址.单",
	AFLD_D:        "浮装载.双",
	AFLD_S:        "浮装载.单",
	AFLOGB_D:      "浮对数基数.双",
	AFLOGB_S:      "浮对数基数.单",
	AFMADD_D:      "浮乘加.双",
	AFMADD_S:      "浮乘加.单",
	AFMAXA_D:      "浮最大绝对值.双",
	AFMAXA_S:      "浮最大绝对值.单",
	AFMAX_D:       "浮最大值.双",
	AFMAX_S:       "浮最大值.单",
	AFMINA_D:      "浮最小绝对值.双",
	AFMINA_S:      "浮最小绝对值.单",
	AFMIN_D:       "浮最小值.双",
	AFMIN_S:       "浮最小值.单",
	AFMOV_D:       "浮移动.双",
	AFMOV_S:       "浮移动.单",
	AFMSUB_D:      "浮乘减.双",
	AFMSUB_S:      "浮乘减.单",
	AFMUL_D:       "浮乘.双",
	AFMUL_S:       "浮乘.单",
	AFNEG_D:       "浮取反.双",
	AFNEG_S:       "浮取反.单",
	AFNMADD_D:     "浮负乘加.双",
	AFNMADD_S:     "浮负乘加.单",
	AFNMSUB_D:     "浮负乘减.双",
	AFNMSUB_S:     "浮负乘减.单",
	AFRECIPE_D:    "浮倒数估值.双",
	AFRECIPE_S:    "浮倒数估值.单",
	AFRECIP_D:     "浮倒数.双",
	AFRECIP_S:     "浮倒数.单",
	AFRINT_D:      "浮舍入到整数.双",
	AFRINT_S:      "浮舍入到整数.单",
	AFRSQRTE_D:    "浮倒数平方根估值.双",
	AFRSQRTE_S:    "浮倒数平方根估值.单",
	AFRSQRT_D:     "浮倒数平方根.双",
	AFRSQRT_S:     "浮倒数平方根.单",
	AFSCALEB_D:    "浮比例基数.双",
	AFSCALEB_S:    "浮比例基数.单",
	AFSEL:         "浮选择",
	AFSQRT_D:      "浮平方根.双",
	AFSQRT_S:      "浮平方根.单",
	AFSTGT_D:      "浮存储.大于.双",
	AFSTGT_S:      "浮存储.大于.单",
	AFSTLE_D:      "浮存储.小于等于.双",
	AFSTLE_S:      "浮存储.小于等于.单",
	AFSTX_D:       "浮存储变址.双",
	AFSTX_S:       "浮存储变址.单",
	AFST_D:        "浮存储.双",
	AFST_S:        "浮存储.单",
	AFSUB_D:       "浮减.双",
	AFSUB_S:       "浮减.单",
	AFTINTRM_L_D:  "浮转定点整数.向负无穷舍入.长.双",
	AFTINTRM_L_S:  "浮转定点整数.向负无穷舍入.长.单",
	AFTINTRM_W_D:  "浮转定点整数.向负无穷舍入.字.双",
	AFTINTRM_W_S:  "浮转定点整数.向负无穷舍入.字.单",
	AFTINTRNE_L_D: "浮转定点整数.向最近舍入.长.双",
	AFTINTRNE_L_S: "浮转定点整数.向最近舍入.长.单",
	AFTINTRNE_W_D: "浮转定点整数.向最近舍入.字.双",
	AFTINTRNE_W_S: "浮转定点整数.向最近舍入.字.单",
	AFTINTRP_L_D:  "浮转定点整数.向正无穷舍入.长.双",
	AFTINTRP_L_S:  "浮转定点整数.向正无穷舍入.长.单",
	AFTINTRP_W_D:  "浮转定点整数.向正无穷舍入.字.双",
	AFTINTRP_W_S:  "浮转定点整数.向正无穷舍入.字.单",
	AFTINTRZ_L_D:  "浮转定点整数.向零截断.长.双",
	AFTINTRZ_L_S:  "浮转定点整数.向零截断.长.单",
	AFTINTRZ_W_D:  "浮转定点整数.向零截断.字.双",
	AFTINTRZ_W_S:  "浮转定点整数.向零截断.字.单",
	AFTINT_L_D:    "浮转定点整数.长.双",
	AFTINT_L_S:    "浮转定点整数.长.单",
	AFTINT_W_D:    "浮转定点整数.字.双",
	AFTINT_W_S:    "浮转定点整数.字.单",
	AIBAR:         "指令屏障",
	AIDLE:         "闲置",
	AINVTLB:       "TLB失效",
	AIOCSRRD_B:    "读IO控制状态.微",
	AIOCSRRD_D:    "读IO控制状态.长",
	AIOCSRRD_H:    "读IO控制状态.短",
	AIOCSRRD_W:    "读IO控制状态.字",
	AIOCSRWR_B:    "写IO控制状态.微",
	AIOCSRWR_D:    "写IO控制状态.长",
	AIOCSRWR_H:    "写IO控制状态.短",
	AIOCSRWR_W:    "写IO控制状态.字",
	AJIRL:         "链接跳转",
	ALDDIR:        "装载直接",
	ALDGT_B:       "装载.大于.微",
	ALDGT_D:       "装载.大于.长",
	ALDGT_H:       "装载.大于.短",
	ALDGT_W:       "装载.大于.字",
	ALDLE_B:       "装载.小于等于.微",
	ALDLE_D:       "装载.小于等于.长",
	ALDLE_H:       "装载.小于等于.短",
	ALDLE_W:       "装载.小于等于.字",
	ALDPTE:        "装载页表项",
	ALDPTR_D:      "装载指针.长",
	ALDPTR_W:      "装载指针.字",
	ALDX_B:        "装载变址.微",
	ALDX_BU:       "装载变址.微正",
	ALDX_D:        "装载变址.长",
	ALDX_H:        "装载变址.短",
	ALDX_HU:       "装载变址.短正",
	ALDX_W:        "装载变址.字",
	ALDX_WU:       "装载变址.字正",
	ALD_B:         "装载.微",
	ALD_BU:        "装载.微正",
	ALD_D:         "装载.长",
	ALD_H:         "装载.短",
	ALD_HU:        "装载.短正",
	ALD_W:         "装载.字",
	ALD_WU:        "装载.字正",
	ALLACQ_D:      "链接装载获取.长",
	ALLACQ_W:      "链接装载获取.字",
	ALL_D:         "链接装载.长",
	ALL_W:         "链接装载.字",
	ALU12I_W:      "装载上部12立.字",
	ALU32I_D:      "装载上部32立.长",
	ALU52I_D:      "装载上部52立.长",
	AMASKEQZ:      "零相等掩码",
	AMASKNEZ:      "零不相等掩码",
	AMOD_D:        "模.长",
	AMOD_DU:       "模.长正",
	AMOD_W:        "模.字",
	AMOD_WU:       "模.字正",
	AMOVCF2FR:     "条件浮点寄存器传到浮点寄存器",
	AMOVCF2GR:     "条件浮点寄存器传到通用寄存器",
	AMOVFCSR2GR:   "浮点控制状态寄存器传到通用寄存器",
	AMOVFR2CF:     "浮点寄存器传到条件浮点寄存器",
	AMOVFR2GR_D:   "浮点寄存器传到通用寄存器.双",
	AMOVFR2GR_S:   "浮点寄存器传到通用寄存器.单",
	AMOVFRH2GR_S:  "浮点寄存器高位传到通用寄存器.单",
	AMOVGR2CF:     "通用寄存器传到条件浮点寄存器",
	AMOVGR2FCSR:   "通用寄存器传到浮点控制状态寄存器",
	AMOVGR2FRH_W:  "通用寄存器传到浮点寄存器高位.字",
	AMOVGR2FR_D:   "通用寄存器传到浮点寄存器.双",
	AMOVGR2FR_W:   "通用寄存器传到浮点寄存器.字",
	AMULH_D:       "乘高位.长",
	AMULH_DU:      "乘高位.长正",
	AMULH_W:       "乘高位.字",
	AMULH_WU:      "乘高位.字正",
	AMULW_D_W:     "乘字.长.字",
	AMULW_D_WU:    "乘字.长.字正",
	AMUL_D:        "乘.长",
	AMUL_W:        "乘.字",
	ANOR:          "或非",
	AOR:           "或",
	AORI:          "或立",
	AORN:          "反或",
	APCADDI:       "计加立",
	APCADDU12I:    "计加高12立",
	APCADDU18I:    "计加高18立",
	APCALAU12I:    "计齐加高12立", // PC ALign Add Upper 12 Immediate
	APRELD:        "预取装载",
	APRELDX:       "预取装载变址",
	ARDTIMEH_W:    "读时间高位.字",
	ARDTIMEL_W:    "读时间低位.字",
	ARDTIME_D:     "读时间.长",
	AREVB_2H:      "字节反转.2短",
	AREVB_2W:      "字节反转.2字",
	AREVB_4H:      "字节反转.4短",
	AREVB_D:       "字节反转.长",
	AREVH_2W:      "短字反转.2字",
	AREVH_D:       "短字反转.长",
	AROTRI_D:      "右旋立.长",
	AROTRI_W:      "右旋立.字",
	AROTR_D:       "右旋.长",
	AROTR_W:       "右旋.字",
	ASCREL_D:      "链接存储释放.长",
	ASCREL_W:      "链接存储释放.字",
	ASC_D:         "链接存储.长",
	ASC_Q:         "链接存储.四",
	ASC_W:         "链接存储.字",
	ASLLI_D:       "左逻辑移位立.长",
	ASLLI_W:       "左逻辑移位立.字",
	ASLL_D:        "左逻辑移位.长",
	ASLL_W:        "左逻辑移位.字",
	ASLT:          "置位小于",
	ASLTI:         "置位小于立",
	ASLTU:         "置位小于.正",
	ASLTUI:        "置位小于立.正",
	ASRAI_D:       "右移立.长",
	ASRAI_W:       "右移立.字",
	ASRA_D:        "右移.长",
	ASRA_W:        "右移.字",
	ASRLI_D:       "佑移立.长",
	ASRLI_W:       "佑移立.字",
	ASRL_D:        "佑移.长",
	ASRL_W:        "佑移.字",
	ASTGT_B:       "存储.大于.微",
	ASTGT_D:       "存储.大于.长",
	ASTGT_H:       "存储.大于.短",
	ASTGT_W:       "存储.大于.字",
	ASTLE_B:       "存储.小于等于.微",
	ASTLE_D:       "存储.小于等于.长",
	ASTLE_H:       "存储.小于等于.短",
	ASTLE_W:       "存储.小于等于.字",
	ASTPTR_D:      "存储指针.长",
	ASTPTR_W:      "存储指针.字",
	ASTX_B:        "存储变址.微",
	ASTX_D:        "存储变址.长",
	ASTX_H:        "存储变址.短",
	ASTX_W:        "存储变址.字",
	AST_B:         "存储.微",
	AST_D:         "存储.长",
	AST_H:         "存储.短",
	AST_W:         "存储.字",
	ASUB_D:        "减.长",
	ASUB_W:        "减.字",
	ASYSCALL:      "系统调用",
	ATLBCLR:       "TLB清除",
	ATLBFILL:      "TLB填充",
	ATLBFLUSH:     "TLB刷新",
	ATLBRD:        "TLB读",
	ATLBSRCH:      "TLB搜索",
	ATLBWR:        "TLB写",
	AXOR:          "异或",
	AXORI:         "异或立",
}
