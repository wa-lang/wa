// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// 注意: 此代码是程序生成, 不要手动修改!!!

package main

// 指令编码信息表
var AOpContextTable = [...]OpContextType{
	{mask: 0xffff8000, value: 0x00108000, name: "ADD.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // ADD.D rd, rj, rk
	{mask: 0xffff8000, value: 0x00100000, name: "ADD.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // ADD.W rd, rj, rk
	{mask: 0xffc00000, value: 0x02c00000, name: "ADDI.D", args: instArgs{Arg_rd, Arg_rj, Arg_si12_21_10}},            // ADDI.D rd, rj, si12
	{mask: 0xffc00000, value: 0x02800000, name: "ADDI.W", args: instArgs{Arg_rd, Arg_rj, Arg_si12_21_10}},            // ADDI.W rd, rj, si12
	{mask: 0xfc000000, value: 0x10000000, name: "ADDU16I.D", args: instArgs{Arg_rd, Arg_rj, Arg_si16_25_10}},         // ADDU16I.D rd, rj, si16
	{mask: 0xfffe0000, value: 0x002c0000, name: "ALSL.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk, Arg_sa2_16_15}},     // ALSL.D rd, rj, rk, sa2
	{mask: 0xfffe0000, value: 0x00040000, name: "ALSL.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk, Arg_sa2_16_15}},     // ALSL.W rd, rj, rk, sa2
	{mask: 0xfffe0000, value: 0x00060000, name: "ALSL.WU", args: instArgs{Arg_rd, Arg_rj, Arg_rk, Arg_sa2_16_15}},    // ALSL.WU rd, rj, rk, sa2
	{mask: 0xffff8000, value: 0x385d0000, name: "AMADD.B", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                   // AMADD.B rd, rk, rj
	{mask: 0xffff8000, value: 0x38618000, name: "AMADD.D", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                   // AMADD.D rd, rk, rj
	{mask: 0xffff8000, value: 0x385d8000, name: "AMADD.H", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                   // AMADD.H rd, rk, rj
	{mask: 0xffff8000, value: 0x38610000, name: "AMADD.W", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                   // AMADD.W rd, rk, rj
	{mask: 0xffff8000, value: 0x385f0000, name: "AMADD_DB.B", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                // AMADD_DB.B rd, rk, rj
	{mask: 0xffff8000, value: 0x386a8000, name: "AMADD_DB.D", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                // AMADD_DB.D rd, rk, rj
	{mask: 0xffff8000, value: 0x385f8000, name: "AMADD_DB.H", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                // AMADD_DB.H rd, rk, rj
	{mask: 0xffff8000, value: 0x386a0000, name: "AMADD_DB.W", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                // AMADD_DB.W rd, rk, rj
	{mask: 0xffff8000, value: 0x38628000, name: "AMAND.D", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                   // AMAND.D rd, rk, rj
	{mask: 0xffff8000, value: 0x38620000, name: "AMAND.W", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                   // AMAND.W rd, rk, rj
	{mask: 0xffff8000, value: 0x386b8000, name: "AMAND_DB.D", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                // AMAND_DB.D rd, rk, rj
	{mask: 0xffff8000, value: 0x386b0000, name: "AMAND_DB.W", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                // AMAND_DB.W rd, rk, rj
	{mask: 0xffff8000, value: 0x38580000, name: "AMCAS.B", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                   // AMCAS.B rd, rk, rj
	{mask: 0xffff8000, value: 0x38598000, name: "AMCAS.D", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                   // AMCAS.D rd, rk, rj
	{mask: 0xffff8000, value: 0x38588000, name: "AMCAS.H", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                   // AMCAS.H rd, rk, rj
	{mask: 0xffff8000, value: 0x38590000, name: "AMCAS.W", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                   // AMCAS.W rd, rk, rj
	{mask: 0xffff8000, value: 0x385a0000, name: "AMCAS_DB.B", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                // AMCAS_DB.B rd, rk, rj
	{mask: 0xffff8000, value: 0x385b8000, name: "AMCAS_DB.D", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                // AMCAS_DB.D rd, rk, rj
	{mask: 0xffff8000, value: 0x385a8000, name: "AMCAS_DB.H", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                // AMCAS_DB.H rd, rk, rj
	{mask: 0xffff8000, value: 0x385b0000, name: "AMCAS_DB.W", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                // AMCAS_DB.W rd, rk, rj
	{mask: 0xffff8000, value: 0x38658000, name: "AMMAX.D", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                   // AMMAX.D rd, rk, rj
	{mask: 0xffff8000, value: 0x38678000, name: "AMMAX.DU", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                  // AMMAX.DU rd, rk, rj
	{mask: 0xffff8000, value: 0x38650000, name: "AMMAX.W", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                   // AMMAX.W rd, rk, rj
	{mask: 0xffff8000, value: 0x38670000, name: "AMMAX.WU", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                  // AMMAX.WU rd, rk, rj
	{mask: 0xffff8000, value: 0x386e8000, name: "AMMAX_DB.D", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                // AMMAX_DB.D rd, rk, rj
	{mask: 0xffff8000, value: 0x38708000, name: "AMMAX_DB.DU", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},               // AMMAX_DB.DU rd, rk, rj
	{mask: 0xffff8000, value: 0x386e0000, name: "AMMAX_DB.W", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                // AMMAX_DB.W rd, rk, rj
	{mask: 0xffff8000, value: 0x38700000, name: "AMMAX_DB.WU", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},               // AMMAX_DB.WU rd, rk, rj
	{mask: 0xffff8000, value: 0x38668000, name: "AMMIN.D", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                   // AMMIN.D rd, rk, rj
	{mask: 0xffff8000, value: 0x38688000, name: "AMMIN.DU", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                  // AMMIN.DU rd, rk, rj
	{mask: 0xffff8000, value: 0x38660000, name: "AMMIN.W", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                   // AMMIN.W rd, rk, rj
	{mask: 0xffff8000, value: 0x38680000, name: "AMMIN.WU", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                  // AMMIN.WU rd, rk, rj
	{mask: 0xffff8000, value: 0x386f8000, name: "AMMIN_DB.D", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                // AMMIN_DB.D rd, rk, rj
	{mask: 0xffff8000, value: 0x38718000, name: "AMMIN_DB.DU", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},               // AMMIN_DB.DU rd, rk, rj
	{mask: 0xffff8000, value: 0x386f0000, name: "AMMIN_DB.W", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                // AMMIN_DB.W rd, rk, rj
	{mask: 0xffff8000, value: 0x38710000, name: "AMMIN_DB.WU", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},               // AMMIN_DB.WU rd, rk, rj
	{mask: 0xffff8000, value: 0x38638000, name: "AMOR.D", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                    // AMOR.D rd, rk, rj
	{mask: 0xffff8000, value: 0x38630000, name: "AMOR.W", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                    // AMOR.W rd, rk, rj
	{mask: 0xffff8000, value: 0x386c8000, name: "AMOR_DB.D", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                 // AMOR_DB.D rd, rk, rj
	{mask: 0xffff8000, value: 0x386c0000, name: "AMOR_DB.W", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                 // AMOR_DB.W rd, rk, rj
	{mask: 0xffff8000, value: 0x385c0000, name: "AMSWAP.B", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                  // AMSWAP.B rd, rk, rj
	{mask: 0xffff8000, value: 0x38608000, name: "AMSWAP.D", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                  // AMSWAP.D rd, rk, rj
	{mask: 0xffff8000, value: 0x385c8000, name: "AMSWAP.H", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                  // AMSWAP.H rd, rk, rj
	{mask: 0xffff8000, value: 0x38600000, name: "AMSWAP.W", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                  // AMSWAP.W rd, rk, rj
	{mask: 0xffff8000, value: 0x385e0000, name: "AMSWAP_DB.B", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},               // AMSWAP_DB.B rd, rk, rj
	{mask: 0xffff8000, value: 0x38698000, name: "AMSWAP_DB.D", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},               // AMSWAP_DB.D rd, rk, rj
	{mask: 0xffff8000, value: 0x385e8000, name: "AMSWAP_DB.H", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},               // AMSWAP_DB.H rd, rk, rj
	{mask: 0xffff8000, value: 0x38690000, name: "AMSWAP_DB.W", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},               // AMSWAP_DB.W rd, rk, rj
	{mask: 0xffff8000, value: 0x38648000, name: "AMXOR.D", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                   // AMXOR.D rd, rk, rj
	{mask: 0xffff8000, value: 0x38640000, name: "AMXOR.W", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                   // AMXOR.W rd, rk, rj
	{mask: 0xffff8000, value: 0x386d8000, name: "AMXOR_DB.D", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                // AMXOR_DB.D rd, rk, rj
	{mask: 0xffff8000, value: 0x386d0000, name: "AMXOR_DB.W", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                // AMXOR_DB.W rd, rk, rj
	{mask: 0xffff8000, value: 0x00148000, name: "AND", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                       // AND rd, rj, rk
	{mask: 0xffc00000, value: 0x03400000, name: "ANDI", args: instArgs{Arg_rd, Arg_rj, Arg_ui12_21_10}},              // ANDI rd, rj, ui12
	{mask: 0xffff8000, value: 0x00168000, name: "ANDN", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                      // ANDN rd, rj, rk
	{mask: 0xffff801f, value: 0x00018000, name: "ASRTGT.D", args: instArgs{Arg_rj, Arg_rk}},                          // ASRTGT.D rj, rk
	{mask: 0xffff801f, value: 0x00010000, name: "ASRTLE.D", args: instArgs{Arg_rj, Arg_rk}},                          // ASRTLE.D rj, rk
	{mask: 0xfc000000, value: 0x50000000, name: "B", args: instArgs{Arg_offset_25_0}},                                // B offs
	{mask: 0xfc000300, value: 0x48000000, name: "BCEQZ", args: instArgs{Arg_cj, Arg_offset_20_0}},                    // BCEQZ cj, offs
	{mask: 0xfc000300, value: 0x48000100, name: "BCNEZ", args: instArgs{Arg_cj, Arg_offset_20_0}},                    // BCNEZ cj, offs
	{mask: 0xfc000000, value: 0x58000000, name: "BEQ", args: instArgs{Arg_rj, Arg_rd, Arg_offset_15_0}},              // BEQ rj, rd, offs
	{mask: 0xfc000000, value: 0x40000000, name: "BEQZ", args: instArgs{Arg_rj, Arg_offset_20_0}},                     // BEQZ rj, offs
	{mask: 0xfc000000, value: 0x64000000, name: "BGE", args: instArgs{Arg_rj, Arg_rd, Arg_offset_15_0}},              // BGE rj, rd, offs
	{mask: 0xfc000000, value: 0x6c000000, name: "BGEU", args: instArgs{Arg_rj, Arg_rd, Arg_offset_15_0}},             // BGEU rj, rd, offs
	{mask: 0xfffffc00, value: 0x00004800, name: "BITREV.4B", args: instArgs{Arg_rd, Arg_rj}},                         // BITREV.4B rd, rj
	{mask: 0xfffffc00, value: 0x00004c00, name: "BITREV.8B", args: instArgs{Arg_rd, Arg_rj}},                         // BITREV.8B rd, rj
	{mask: 0xfffffc00, value: 0x00005400, name: "BITREV.D", args: instArgs{Arg_rd, Arg_rj}},                          // BITREV.D rd, rj
	{mask: 0xfffffc00, value: 0x00005000, name: "BITREV.W", args: instArgs{Arg_rd, Arg_rj}},                          // BITREV.W rd, rj
	{mask: 0xfc000000, value: 0x54000000, name: "BL", args: instArgs{Arg_offset_25_0}},                               // BL offs
	{mask: 0xfc000000, value: 0x60000000, name: "BLT", args: instArgs{Arg_rj, Arg_rd, Arg_offset_15_0}},              // BLT rj, rd, offs
	{mask: 0xfc000000, value: 0x68000000, name: "BLTU", args: instArgs{Arg_rj, Arg_rd, Arg_offset_15_0}},             // BLTU rj, rd, offs
	{mask: 0xfc000000, value: 0x5c000000, name: "BNE", args: instArgs{Arg_rj, Arg_rd, Arg_offset_15_0}},              // BNE rj, rd, offs
	{mask: 0xfc000000, value: 0x44000000, name: "BNEZ", args: instArgs{Arg_rj, Arg_offset_20_0}},                     // BNEZ rj, offs
	{mask: 0xffff8000, value: 0x002a0000, name: "BREAK", args: instArgs{Arg_code_14_0}},                              // BREAK code
	{mask: 0xffc00000, value: 0x00800000, name: "BSTRINS.D", args: instArgs{Arg_rd, Arg_rj, Arg_msbd, Arg_lsbd}},     // BSTRINS.D rd, rj, msbd, lsbd
	{mask: 0xffe08000, value: 0x00600000, name: "BSTRINS.W", args: instArgs{Arg_rd, Arg_rj, Arg_msbw, Arg_lsbw}},     // BSTRINS.W rd, rj, msbw, lsbw
	{mask: 0xffc00000, value: 0x00c00000, name: "BSTRPICK.D", args: instArgs{Arg_rd, Arg_rj, Arg_msbd, Arg_lsbd}},    // BSTRPICK.D rd, rj, msbd, lsbd
	{mask: 0xffe08000, value: 0x00608000, name: "BSTRPICK.W", args: instArgs{Arg_rd, Arg_rj, Arg_msbw, Arg_lsbw}},    // BSTRPICK.W rd, rj, msbw, lsbw
	{mask: 0xfffc0000, value: 0x000c0000, name: "BYTEPICK.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk, Arg_sa3_17_15}}, // BYTEPICK.D rd, rj, rk, sa3
	{mask: 0xfffe0000, value: 0x00080000, name: "BYTEPICK.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk, Arg_sa2_16_15}}, // BYTEPICK.W rd, rj, rk, sa2
	{mask: 0xffc00000, value: 0x06000000, name: "CACOP", args: instArgs{Arg_code_4_0, Arg_rj, Arg_si12_21_10}},       // CACOP code, rj, si12
	{mask: 0xfffffc00, value: 0x00002000, name: "CLO.D", args: instArgs{Arg_rd, Arg_rj}},                             // CLO.D rd, rj
	{mask: 0xfffffc00, value: 0x00001000, name: "CLO.W", args: instArgs{Arg_rd, Arg_rj}},                             // CLO.W rd, rj
	{mask: 0xfffffc00, value: 0x00002400, name: "CLZ.D", args: instArgs{Arg_rd, Arg_rj}},                             // CLZ.D rd, rj
	{mask: 0xfffffc00, value: 0x00001400, name: "CLZ.W", args: instArgs{Arg_rd, Arg_rj}},                             // CLZ.W rd, rj
	{mask: 0xfffffc00, value: 0x00006c00, name: "CPUCFG", args: instArgs{Arg_rd, Arg_rj}},                            // CPUCFG rd, rj
	{mask: 0xffff8000, value: 0x00240000, name: "CRC.W.B.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                 // CRC.W.B.W rd, rj, rk
	{mask: 0xffff8000, value: 0x00258000, name: "CRC.W.D.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                 // CRC.W.D.W rd, rj, rk
	{mask: 0xffff8000, value: 0x00248000, name: "CRC.W.H.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                 // CRC.W.H.W rd, rj, rk
	{mask: 0xffff8000, value: 0x00250000, name: "CRC.W.W.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                 // CRC.W.W.W rd, rj, rk
	{mask: 0xffff8000, value: 0x00260000, name: "CRCC.W.B.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                // CRCC.W.B.W rd, rj, rk
	{mask: 0xffff8000, value: 0x00278000, name: "CRCC.W.D.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                // CRCC.W.D.W rd, rj, rk
	{mask: 0xffff8000, value: 0x00268000, name: "CRCC.W.H.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                // CRCC.W.H.W rd, rj, rk
	{mask: 0xffff8000, value: 0x00270000, name: "CRCC.W.W.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                // CRCC.W.W.W rd, rj, rk
	{mask: 0xff0003e0, value: 0x04000000, name: "CSRRD", args: instArgs{Arg_rd, Arg_csr_23_10}},                      // CSRRD rd, csr
	{mask: 0xff0003e0, value: 0x04000020, name: "CSRWR", args: instArgs{Arg_rd, Arg_csr_23_10}},                      // CSRWR rd, csr
	{mask: 0xff000000, value: 0x04000000, name: "CSRXCHG", args: instArgs{Arg_rd, Arg_rj, Arg_csr_23_10}},            // CSRXCHG rd, rj, csr
	{mask: 0xfffffc00, value: 0x00002800, name: "CTO.D", args: instArgs{Arg_rd, Arg_rj}},                             // CTO.D rd, rj
	{mask: 0xfffffc00, value: 0x00001800, name: "CTO.W", args: instArgs{Arg_rd, Arg_rj}},                             // CTO.W rd, rj
	{mask: 0xfffffc00, value: 0x00002c00, name: "CTZ.D", args: instArgs{Arg_rd, Arg_rj}},                             // CTZ.D rd, rj
	{mask: 0xfffffc00, value: 0x00001c00, name: "CTZ.W", args: instArgs{Arg_rd, Arg_rj}},                             // CTZ.W rd, rj
	{mask: 0xffff8000, value: 0x38720000, name: "DBAR", args: instArgs{Arg_hint_14_0}},                               // DBAR hint
	{mask: 0xffff8000, value: 0x002a8000, name: "DBCL", args: instArgs{Arg_code_14_0}},                               // DBCL code
	{mask: 0xffff8000, value: 0x00220000, name: "DIV.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // DIV.D rd, rj, rk
	{mask: 0xffff8000, value: 0x00230000, name: "DIV.DU", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // DIV.DU rd, rj, rk
	{mask: 0xffff8000, value: 0x00200000, name: "DIV.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // DIV.W rd, rj, rk
	{mask: 0xffff8000, value: 0x00210000, name: "DIV.WU", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // DIV.WU rd, rj, rk
	{mask: 0xffffffff, value: 0x06483800, name: "ERTN", args: instArgs{}},                                            // ERTN
	{mask: 0xfffffc00, value: 0x00005c00, name: "EXT.W.B", args: instArgs{Arg_rd, Arg_rj}},                           // EXT.W.B rd, rj
	{mask: 0xfffffc00, value: 0x00005800, name: "EXT.W.H", args: instArgs{Arg_rd, Arg_rj}},                           // EXT.W.H rd, rj
	{mask: 0xfffffc00, value: 0x01140800, name: "FABS.D", args: instArgs{Arg_fd, Arg_fj}},                            // FABS.D fd, fj
	{mask: 0xfffffc00, value: 0x01140400, name: "FABS.S", args: instArgs{Arg_fd, Arg_fj}},                            // FABS.S fd, fj
	{mask: 0xffff8000, value: 0x01010000, name: "FADD.D", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                    // FADD.D fd, fj, fk
	{mask: 0xffff8000, value: 0x01008000, name: "FADD.S", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                    // FADD.S fd, fj, fk
	{mask: 0xfffffc00, value: 0x01143800, name: "FCLASS.D", args: instArgs{Arg_fd, Arg_fj}},                          // FCLASS.D fd, fj
	{mask: 0xfffffc00, value: 0x01143400, name: "FCLASS.S", args: instArgs{Arg_fd, Arg_fj}},                          // FCLASS.S fd, fj
	{mask: 0xffff8018, value: 0x0c200000, name: "FCMP.CAF.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.CAF.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c100000, name: "FCMP.CAF.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.CAF.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c220000, name: "FCMP.CEQ.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.CEQ.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c120000, name: "FCMP.CEQ.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.CEQ.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c230000, name: "FCMP.CLE.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.CLE.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c130000, name: "FCMP.CLE.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.CLE.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c210000, name: "FCMP.CLT.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.CLT.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c110000, name: "FCMP.CLT.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.CLT.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c280000, name: "FCMP.CNE.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.CNE.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c180000, name: "FCMP.CNE.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.CNE.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c2a0000, name: "FCMP.COR.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.COR.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c1a0000, name: "FCMP.COR.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.COR.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c260000, name: "FCMP.CUEQ.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},               // FCMP.CUEQ.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c160000, name: "FCMP.CUEQ.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},               // FCMP.CUEQ.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c270000, name: "FCMP.CULE.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},               // FCMP.CULE.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c170000, name: "FCMP.CULE.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},               // FCMP.CULE.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c250000, name: "FCMP.CULT.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},               // FCMP.CULT.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c150000, name: "FCMP.CULT.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},               // FCMP.CULT.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c240000, name: "FCMP.CUN.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.CUN.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c140000, name: "FCMP.CUN.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.CUN.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c2c0000, name: "FCMP.CUNE.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},               // FCMP.CUNE.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c1c0000, name: "FCMP.CUNE.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},               // FCMP.CUNE.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c208000, name: "FCMP.SAF.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.SAF.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c108000, name: "FCMP.SAF.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.SAF.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c228000, name: "FCMP.SEQ.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.SEQ.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c128000, name: "FCMP.SEQ.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.SEQ.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c238000, name: "FCMP.SLE.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.SLE.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c138000, name: "FCMP.SLE.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.SLE.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c218000, name: "FCMP.SLT.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.SLT.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c118000, name: "FCMP.SLT.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.SLT.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c288000, name: "FCMP.SNE.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.SNE.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c188000, name: "FCMP.SNE.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.SNE.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c2a8000, name: "FCMP.SOR.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.SOR.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c1a8000, name: "FCMP.SOR.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.SOR.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c268000, name: "FCMP.SUEQ.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},               // FCMP.SUEQ.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c168000, name: "FCMP.SUEQ.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},               // FCMP.SUEQ.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c278000, name: "FCMP.SULE.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},               // FCMP.SULE.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c178000, name: "FCMP.SULE.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},               // FCMP.SULE.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c258000, name: "FCMP.SULT.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},               // FCMP.SULT.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c158000, name: "FCMP.SULT.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},               // FCMP.SULT.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c248000, name: "FCMP.SUN.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.SUN.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c148000, name: "FCMP.SUN.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},                // FCMP.SUN.S cd, fj, fk
	{mask: 0xffff8018, value: 0x0c2c8000, name: "FCMP.SUNE.D", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},               // FCMP.SUNE.D cd, fj, fk
	{mask: 0xffff8018, value: 0x0c1c8000, name: "FCMP.SUNE.S", args: instArgs{Arg_cd, Arg_fj, Arg_fk}},               // FCMP.SUNE.S cd, fj, fk
	{mask: 0xffff8000, value: 0x01130000, name: "FCOPYSIGN.D", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},               // FCOPYSIGN.D fd, fj, fk
	{mask: 0xffff8000, value: 0x01128000, name: "FCOPYSIGN.S", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},               // FCOPYSIGN.S fd, fj, fk
	{mask: 0xfffffc00, value: 0x01192400, name: "FCVT.D.S", args: instArgs{Arg_fd, Arg_fj}},                          // FCVT.D.S fd, fj
	{mask: 0xfffffc00, value: 0x01191800, name: "FCVT.S.D", args: instArgs{Arg_fd, Arg_fj}},                          // FCVT.S.D fd, fj
	{mask: 0xffff8000, value: 0x01070000, name: "FDIV.D", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                    // FDIV.D fd, fj, fk
	{mask: 0xffff8000, value: 0x01068000, name: "FDIV.S", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                    // FDIV.S fd, fj, fk
	{mask: 0xfffffc00, value: 0x011d2800, name: "FFINT.D.L", args: instArgs{Arg_fd, Arg_fj}},                         // FFINT.D.L fd, fj
	{mask: 0xfffffc00, value: 0x011d2000, name: "FFINT.D.W", args: instArgs{Arg_fd, Arg_fj}},                         // FFINT.D.W fd, fj
	{mask: 0xfffffc00, value: 0x011d1800, name: "FFINT.S.L", args: instArgs{Arg_fd, Arg_fj}},                         // FFINT.S.L fd, fj
	{mask: 0xfffffc00, value: 0x011d1000, name: "FFINT.S.W", args: instArgs{Arg_fd, Arg_fj}},                         // FFINT.S.W fd, fj
	{mask: 0xffc00000, value: 0x2b800000, name: "FLD.D", args: instArgs{Arg_fd, Arg_rj, Arg_si12_21_10}},             // FLD.D fd, rj, si12
	{mask: 0xffc00000, value: 0x2b000000, name: "FLD.S", args: instArgs{Arg_fd, Arg_rj, Arg_si12_21_10}},             // FLD.S fd, rj, si12
	{mask: 0xffff8000, value: 0x38748000, name: "FLDGT.D", args: instArgs{Arg_fd, Arg_rj, Arg_rk}},                   // FLDGT.D fd, rj, rk
	{mask: 0xffff8000, value: 0x38740000, name: "FLDGT.S", args: instArgs{Arg_fd, Arg_rj, Arg_rk}},                   // FLDGT.S fd, rj, rk
	{mask: 0xffff8000, value: 0x38758000, name: "FLDLE.D", args: instArgs{Arg_fd, Arg_rj, Arg_rk}},                   // FLDLE.D fd, rj, rk
	{mask: 0xffff8000, value: 0x38750000, name: "FLDLE.S", args: instArgs{Arg_fd, Arg_rj, Arg_rk}},                   // FLDLE.S fd, rj, rk
	{mask: 0xffff8000, value: 0x38340000, name: "FLDX.D", args: instArgs{Arg_fd, Arg_rj, Arg_rk}},                    // FLDX.D fd, rj, rk
	{mask: 0xffff8000, value: 0x38300000, name: "FLDX.S", args: instArgs{Arg_fd, Arg_rj, Arg_rk}},                    // FLDX.S fd, rj, rk
	{mask: 0xfffffc00, value: 0x01142800, name: "FLOGB.D", args: instArgs{Arg_fd, Arg_fj}},                           // FLOGB.D fd, fj
	{mask: 0xfffffc00, value: 0x01142400, name: "FLOGB.S", args: instArgs{Arg_fd, Arg_fj}},                           // FLOGB.S fd, fj
	{mask: 0xfff00000, value: 0x08200000, name: "FMADD.D", args: instArgs{Arg_fd, Arg_fj, Arg_fk, Arg_fa}},           // FMADD.D fd, fj, fk, fa
	{mask: 0xfff00000, value: 0x08100000, name: "FMADD.S", args: instArgs{Arg_fd, Arg_fj, Arg_fk, Arg_fa}},           // FMADD.S fd, fj, fk, fa
	{mask: 0xffff8000, value: 0x01090000, name: "FMAX.D", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                    // FMAX.D fd, fj, fk
	{mask: 0xffff8000, value: 0x01088000, name: "FMAX.S", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                    // FMAX.S fd, fj, fk
	{mask: 0xffff8000, value: 0x010d0000, name: "FMAXA.D", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                   // FMAXA.D fd, fj, fk
	{mask: 0xffff8000, value: 0x010c8000, name: "FMAXA.S", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                   // FMAXA.S fd, fj, fk
	{mask: 0xffff8000, value: 0x010b0000, name: "FMIN.D", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                    // FMIN.D fd, fj, fk
	{mask: 0xffff8000, value: 0x010a8000, name: "FMIN.S", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                    // FMIN.S fd, fj, fk
	{mask: 0xffff8000, value: 0x010f0000, name: "FMINA.D", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                   // FMINA.D fd, fj, fk
	{mask: 0xffff8000, value: 0x010e8000, name: "FMINA.S", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                   // FMINA.S fd, fj, fk
	{mask: 0xfffffc00, value: 0x01149800, name: "FMOV.D", args: instArgs{Arg_fd, Arg_fj}},                            // FMOV.D fd, fj
	{mask: 0xfffffc00, value: 0x01149400, name: "FMOV.S", args: instArgs{Arg_fd, Arg_fj}},                            // FMOV.S fd, fj
	{mask: 0xfff00000, value: 0x08600000, name: "FMSUB.D", args: instArgs{Arg_fd, Arg_fj, Arg_fk, Arg_fa}},           // FMSUB.D fd, fj, fk, fa
	{mask: 0xfff00000, value: 0x08500000, name: "FMSUB.S", args: instArgs{Arg_fd, Arg_fj, Arg_fk, Arg_fa}},           // FMSUB.S fd, fj, fk, fa
	{mask: 0xffff8000, value: 0x01050000, name: "FMUL.D", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                    // FMUL.D fd, fj, fk
	{mask: 0xffff8000, value: 0x01048000, name: "FMUL.S", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                    // FMUL.S fd, fj, fk
	{mask: 0xfffffc00, value: 0x01141800, name: "FNEG.D", args: instArgs{Arg_fd, Arg_fj}},                            // FNEG.D fd, fj
	{mask: 0xfffffc00, value: 0x01141400, name: "FNEG.S", args: instArgs{Arg_fd, Arg_fj}},                            // FNEG.S fd, fj
	{mask: 0xfff00000, value: 0x08a00000, name: "FNMADD.D", args: instArgs{Arg_fd, Arg_fj, Arg_fk, Arg_fa}},          // FNMADD.D fd, fj, fk, fa
	{mask: 0xfff00000, value: 0x08900000, name: "FNMADD.S", args: instArgs{Arg_fd, Arg_fj, Arg_fk, Arg_fa}},          // FNMADD.S fd, fj, fk, fa
	{mask: 0xfff00000, value: 0x08e00000, name: "FNMSUB.D", args: instArgs{Arg_fd, Arg_fj, Arg_fk, Arg_fa}},          // FNMSUB.D fd, fj, fk, fa
	{mask: 0xfff00000, value: 0x08d00000, name: "FNMSUB.S", args: instArgs{Arg_fd, Arg_fj, Arg_fk, Arg_fa}},          // FNMSUB.S fd, fj, fk, fa
	{mask: 0xfffffc00, value: 0x01145800, name: "FRECIP.D", args: instArgs{Arg_fd, Arg_fj}},                          // FRECIP.D fd, fj
	{mask: 0xfffffc00, value: 0x01145400, name: "FRECIP.S", args: instArgs{Arg_fd, Arg_fj}},                          // FRECIP.S fd, fj
	{mask: 0xfffffc00, value: 0x01147800, name: "FRECIPE.D", args: instArgs{Arg_fd, Arg_fj}},                         // FRECIPE.D fd, fj
	{mask: 0xfffffc00, value: 0x01147400, name: "FRECIPE.S", args: instArgs{Arg_fd, Arg_fj}},                         // FRECIPE.S fd, fj
	{mask: 0xfffffc00, value: 0x011e4800, name: "FRINT.D", args: instArgs{Arg_fd, Arg_fj}},                           // FRINT.D fd, fj
	{mask: 0xfffffc00, value: 0x011e4400, name: "FRINT.S", args: instArgs{Arg_fd, Arg_fj}},                           // FRINT.S fd, fj
	{mask: 0xfffffc00, value: 0x01146800, name: "FRSQRT.D", args: instArgs{Arg_fd, Arg_fj}},                          // FRSQRT.D fd, fj
	{mask: 0xfffffc00, value: 0x01146400, name: "FRSQRT.S", args: instArgs{Arg_fd, Arg_fj}},                          // FRSQRT.S fd, fj
	{mask: 0xfffffc00, value: 0x01148800, name: "FRSQRTE.D", args: instArgs{Arg_fd, Arg_fj}},                         // FRSQRTE.D fd, fj
	{mask: 0xfffffc00, value: 0x01148400, name: "FRSQRTE.S", args: instArgs{Arg_fd, Arg_fj}},                         // FRSQRTE.S fd, fj
	{mask: 0xffff8000, value: 0x01110000, name: "FSCALEB.D", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                 // FSCALEB.D fd, fj, fk
	{mask: 0xffff8000, value: 0x01108000, name: "FSCALEB.S", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                 // FSCALEB.S fd, fj, fk
	{mask: 0xfffc0000, value: 0x0d000000, name: "FSEL", args: instArgs{Arg_fd, Arg_fj, Arg_fk, Arg_ca}},              // FSEL fd, fj, fk, ca
	{mask: 0xfffffc00, value: 0x01144800, name: "FSQRT.D", args: instArgs{Arg_fd, Arg_fj}},                           // FSQRT.D fd, fj
	{mask: 0xfffffc00, value: 0x01144400, name: "FSQRT.S", args: instArgs{Arg_fd, Arg_fj}},                           // FSQRT.S fd, fj
	{mask: 0xffc00000, value: 0x2bc00000, name: "FST.D", args: instArgs{Arg_fd, Arg_rj, Arg_si12_21_10}},             // FST.D fd, rj, si12
	{mask: 0xffc00000, value: 0x2b400000, name: "FST.S", args: instArgs{Arg_fd, Arg_rj, Arg_si12_21_10}},             // FST.S fd, rj, si12
	{mask: 0xffff8000, value: 0x38768000, name: "FSTGT.D", args: instArgs{Arg_fd, Arg_rj, Arg_rk}},                   // FSTGT.D fd, rj, rk
	{mask: 0xffff8000, value: 0x38760000, name: "FSTGT.S", args: instArgs{Arg_fd, Arg_rj, Arg_rk}},                   // FSTGT.S fd, rj, rk
	{mask: 0xffff8000, value: 0x38778000, name: "FSTLE.D", args: instArgs{Arg_fd, Arg_rj, Arg_rk}},                   // FSTLE.D fd, rj, rk
	{mask: 0xffff8000, value: 0x38770000, name: "FSTLE.S", args: instArgs{Arg_fd, Arg_rj, Arg_rk}},                   // FSTLE.S fd, rj, rk
	{mask: 0xffff8000, value: 0x383c0000, name: "FSTX.D", args: instArgs{Arg_fd, Arg_rj, Arg_rk}},                    // FSTX.D fd, rj, rk
	{mask: 0xffff8000, value: 0x38380000, name: "FSTX.S", args: instArgs{Arg_fd, Arg_rj, Arg_rk}},                    // FSTX.S fd, rj, rk
	{mask: 0xffff8000, value: 0x01030000, name: "FSUB.D", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                    // FSUB.D fd, fj, fk
	{mask: 0xffff8000, value: 0x01028000, name: "FSUB.S", args: instArgs{Arg_fd, Arg_fj, Arg_fk}},                    // FSUB.S fd, fj, fk
	{mask: 0xfffffc00, value: 0x011b2800, name: "FTINT.L.D", args: instArgs{Arg_fd, Arg_fj}},                         // FTINT.L.D fd, fj
	{mask: 0xfffffc00, value: 0x011b2400, name: "FTINT.L.S", args: instArgs{Arg_fd, Arg_fj}},                         // FTINT.L.S fd, fj
	{mask: 0xfffffc00, value: 0x011b0800, name: "FTINT.W.D", args: instArgs{Arg_fd, Arg_fj}},                         // FTINT.W.D fd, fj
	{mask: 0xfffffc00, value: 0x011b0400, name: "FTINT.W.S", args: instArgs{Arg_fd, Arg_fj}},                         // FTINT.W.S fd, fj
	{mask: 0xfffffc00, value: 0x011a2800, name: "FTINTRM.L.D", args: instArgs{Arg_fd, Arg_fj}},                       // FTINTRM.L.D fd, fj
	{mask: 0xfffffc00, value: 0x011a2400, name: "FTINTRM.L.S", args: instArgs{Arg_fd, Arg_fj}},                       // FTINTRM.L.S fd, fj
	{mask: 0xfffffc00, value: 0x011a0800, name: "FTINTRM.W.D", args: instArgs{Arg_fd, Arg_fj}},                       // FTINTRM.W.D fd, fj
	{mask: 0xfffffc00, value: 0x011a0400, name: "FTINTRM.W.S", args: instArgs{Arg_fd, Arg_fj}},                       // FTINTRM.W.S fd, fj
	{mask: 0xfffffc00, value: 0x011ae800, name: "FTINTRNE.L.D", args: instArgs{Arg_fd, Arg_fj}},                      // FTINTRNE.L.D fd, fj
	{mask: 0xfffffc00, value: 0x011ae400, name: "FTINTRNE.L.S", args: instArgs{Arg_fd, Arg_fj}},                      // FTINTRNE.L.S fd, fj
	{mask: 0xfffffc00, value: 0x011ac800, name: "FTINTRNE.W.D", args: instArgs{Arg_fd, Arg_fj}},                      // FTINTRNE.W.D fd, fj
	{mask: 0xfffffc00, value: 0x011ac400, name: "FTINTRNE.W.S", args: instArgs{Arg_fd, Arg_fj}},                      // FTINTRNE.W.S fd, fj
	{mask: 0xfffffc00, value: 0x011a6800, name: "FTINTRP.L.D", args: instArgs{Arg_fd, Arg_fj}},                       // FTINTRP.L.D fd, fj
	{mask: 0xfffffc00, value: 0x011a6400, name: "FTINTRP.L.S", args: instArgs{Arg_fd, Arg_fj}},                       // FTINTRP.L.S fd, fj
	{mask: 0xfffffc00, value: 0x011a4800, name: "FTINTRP.W.D", args: instArgs{Arg_fd, Arg_fj}},                       // FTINTRP.W.D fd, fj
	{mask: 0xfffffc00, value: 0x011a4400, name: "FTINTRP.W.S", args: instArgs{Arg_fd, Arg_fj}},                       // FTINTRP.W.S fd, fj
	{mask: 0xfffffc00, value: 0x011aa800, name: "FTINTRZ.L.D", args: instArgs{Arg_fd, Arg_fj}},                       // FTINTRZ.L.D fd, fj
	{mask: 0xfffffc00, value: 0x011aa400, name: "FTINTRZ.L.S", args: instArgs{Arg_fd, Arg_fj}},                       // FTINTRZ.L.S fd, fj
	{mask: 0xfffffc00, value: 0x011a8800, name: "FTINTRZ.W.D", args: instArgs{Arg_fd, Arg_fj}},                       // FTINTRZ.W.D fd, fj
	{mask: 0xfffffc00, value: 0x011a8400, name: "FTINTRZ.W.S", args: instArgs{Arg_fd, Arg_fj}},                       // FTINTRZ.W.S fd, fj
	{mask: 0xffff8000, value: 0x38728000, name: "IBAR", args: instArgs{Arg_hint_14_0}},                               // IBAR hint
	{mask: 0xffff8000, value: 0x06488000, name: "IDLE", args: instArgs{Arg_level_14_0}},                              // IDLE level
	{mask: 0xffff8000, value: 0x06498000, name: "INVTLB", args: instArgs{Arg_op_4_0, Arg_rj, Arg_rk}},                // INVTLB op, rj, rk
	{mask: 0xfffffc00, value: 0x06480000, name: "IOCSRRD.B", args: instArgs{Arg_rd, Arg_rj}},                         // IOCSRRD.B rd, rj
	{mask: 0xfffffc00, value: 0x06480c00, name: "IOCSRRD.D", args: instArgs{Arg_rd, Arg_rj}},                         // IOCSRRD.D rd, rj
	{mask: 0xfffffc00, value: 0x06480400, name: "IOCSRRD.H", args: instArgs{Arg_rd, Arg_rj}},                         // IOCSRRD.H rd, rj
	{mask: 0xfffffc00, value: 0x06480800, name: "IOCSRRD.W", args: instArgs{Arg_rd, Arg_rj}},                         // IOCSRRD.W rd, rj
	{mask: 0xfffffc00, value: 0x06481000, name: "IOCSRWR.B", args: instArgs{Arg_rd, Arg_rj}},                         // IOCSRWR.B rd, rj
	{mask: 0xfffffc00, value: 0x06481c00, name: "IOCSRWR.D", args: instArgs{Arg_rd, Arg_rj}},                         // IOCSRWR.D rd, rj
	{mask: 0xfffffc00, value: 0x06481400, name: "IOCSRWR.H", args: instArgs{Arg_rd, Arg_rj}},                         // IOCSRWR.H rd, rj
	{mask: 0xfffffc00, value: 0x06481800, name: "IOCSRWR.W", args: instArgs{Arg_rd, Arg_rj}},                         // IOCSRWR.W rd, rj
	{mask: 0xfc000000, value: 0x4c000000, name: "JIRL", args: instArgs{Arg_rd, Arg_rj, Arg_offset_15_0}},             // JIRL rd, rj, offs
	{mask: 0xffc00000, value: 0x28000000, name: "LD.B", args: instArgs{Arg_rd, Arg_rj, Arg_si12_21_10}},              // LD.B rd, rj, si12
	{mask: 0xffc00000, value: 0x2a000000, name: "LD.BU", args: instArgs{Arg_rd, Arg_rj, Arg_si12_21_10}},             // LD.BU rd, rj, si12
	{mask: 0xffc00000, value: 0x28c00000, name: "LD.D", args: instArgs{Arg_rd, Arg_rj, Arg_si12_21_10}},              // LD.D rd, rj, si12
	{mask: 0xffc00000, value: 0x28400000, name: "LD.H", args: instArgs{Arg_rd, Arg_rj, Arg_si12_21_10}},              // LD.H rd, rj, si12
	{mask: 0xffc00000, value: 0x2a400000, name: "LD.HU", args: instArgs{Arg_rd, Arg_rj, Arg_si12_21_10}},             // LD.HU rd, rj, si12
	{mask: 0xffc00000, value: 0x28800000, name: "LD.W", args: instArgs{Arg_rd, Arg_rj, Arg_si12_21_10}},              // LD.W rd, rj, si12
	{mask: 0xffc00000, value: 0x2a800000, name: "LD.WU", args: instArgs{Arg_rd, Arg_rj, Arg_si12_21_10}},             // LD.WU rd, rj, si12
	{mask: 0xfffc0000, value: 0x06400000, name: "LDDIR", args: instArgs{Arg_rd, Arg_rj, Arg_level_17_10}},            // LDDIR rd, rj, level
	{mask: 0xffff8000, value: 0x38780000, name: "LDGT.B", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // LDGT.B rd, rj, rk
	{mask: 0xffff8000, value: 0x38798000, name: "LDGT.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // LDGT.D rd, rj, rk
	{mask: 0xffff8000, value: 0x38788000, name: "LDGT.H", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // LDGT.H rd, rj, rk
	{mask: 0xffff8000, value: 0x38790000, name: "LDGT.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // LDGT.W rd, rj, rk
	{mask: 0xffff8000, value: 0x387a0000, name: "LDLE.B", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // LDLE.B rd, rj, rk
	{mask: 0xffff8000, value: 0x387b8000, name: "LDLE.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // LDLE.D rd, rj, rk
	{mask: 0xffff8000, value: 0x387a8000, name: "LDLE.H", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // LDLE.H rd, rj, rk
	{mask: 0xffff8000, value: 0x387b0000, name: "LDLE.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // LDLE.W rd, rj, rk
	{mask: 0xfffc001f, value: 0x06440000, name: "LDPTE", args: instArgs{Arg_rj, Arg_seq_17_10}},                      // LDPTE rj, seq
	{mask: 0xff000000, value: 0x26000000, name: "LDPTR.D", args: instArgs{Arg_rd, Arg_rj, Arg_si14_23_10}},           // LDPTR.D rd, rj, si14
	{mask: 0xff000000, value: 0x24000000, name: "LDPTR.W", args: instArgs{Arg_rd, Arg_rj, Arg_si14_23_10}},           // LDPTR.W rd, rj, si14
	{mask: 0xffff8000, value: 0x38000000, name: "LDX.B", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // LDX.B rd, rj, rk
	{mask: 0xffff8000, value: 0x38200000, name: "LDX.BU", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // LDX.BU rd, rj, rk
	{mask: 0xffff8000, value: 0x380c0000, name: "LDX.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // LDX.D rd, rj, rk
	{mask: 0xffff8000, value: 0x38040000, name: "LDX.H", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // LDX.H rd, rj, rk
	{mask: 0xffff8000, value: 0x38240000, name: "LDX.HU", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // LDX.HU rd, rj, rk
	{mask: 0xffff8000, value: 0x38080000, name: "LDX.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // LDX.W rd, rj, rk
	{mask: 0xffff8000, value: 0x38280000, name: "LDX.WU", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // LDX.WU rd, rj, rk
	{mask: 0xff000000, value: 0x22000000, name: "LL.D", args: instArgs{Arg_rd, Arg_rj, Arg_si14_23_10}},              // LL.D rd, rj, si14
	{mask: 0xff000000, value: 0x20000000, name: "LL.W", args: instArgs{Arg_rd, Arg_rj, Arg_si14_23_10}},              // LL.W rd, rj, si14
	{mask: 0xfffffc00, value: 0x38578800, name: "LLACQ.D", args: instArgs{Arg_rd, Arg_rj}},                           // LLACQ.D rd, rj
	{mask: 0xfffffc00, value: 0x38578000, name: "LLACQ.W", args: instArgs{Arg_rd, Arg_rj}},                           // LLACQ.W rd, rj
	{mask: 0xfe000000, value: 0x14000000, name: "LU12I.W", args: instArgs{Arg_rd, Arg_si20_24_5}},                    // LU12I.W rd, si20
	{mask: 0xfe000000, value: 0x16000000, name: "LU32I.D", args: instArgs{Arg_rd, Arg_si20_24_5}},                    // LU32I.D rd, si20
	{mask: 0xffc00000, value: 0x03000000, name: "LU52I.D", args: instArgs{Arg_rd, Arg_rj, Arg_si12_21_10}},           // LU52I.D rd, rj, si12
	{mask: 0xffff8000, value: 0x00130000, name: "MASKEQZ", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                   // MASKEQZ rd, rj, rk
	{mask: 0xffff8000, value: 0x00138000, name: "MASKNEZ", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                   // MASKNEZ rd, rj, rk
	{mask: 0xffff8000, value: 0x00228000, name: "MOD.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // MOD.D rd, rj, rk
	{mask: 0xffff8000, value: 0x00238000, name: "MOD.DU", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // MOD.DU rd, rj, rk
	{mask: 0xffff8000, value: 0x00208000, name: "MOD.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // MOD.W rd, rj, rk
	{mask: 0xffff8000, value: 0x00218000, name: "MOD.WU", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // MOD.WU rd, rj, rk
	{mask: 0xffffff00, value: 0x0114d400, name: "MOVCF2FR", args: instArgs{Arg_fd, Arg_cj}},                          // MOVCF2FR fd, cj
	{mask: 0xffffff00, value: 0x0114dc00, name: "MOVCF2GR", args: instArgs{Arg_rd, Arg_cj}},                          // MOVCF2GR rd, cj
	{mask: 0xfffffc00, value: 0x0114c800, name: "MOVFCSR2GR", args: instArgs{Arg_rd, Arg_fcsr_9_5}},                  // MOVFCSR2GR rd, fcsr
	{mask: 0xfffffc18, value: 0x0114d000, name: "MOVFR2CF", args: instArgs{Arg_cd, Arg_fj}},                          // MOVFR2CF cd, fj
	{mask: 0xfffffc00, value: 0x0114b800, name: "MOVFR2GR.D", args: instArgs{Arg_rd, Arg_fj}},                        // MOVFR2GR.D rd, fj
	{mask: 0xfffffc00, value: 0x0114b400, name: "MOVFR2GR.S", args: instArgs{Arg_rd, Arg_fj}},                        // MOVFR2GR.S rd, fj
	{mask: 0xfffffc00, value: 0x0114bc00, name: "MOVFRH2GR.S", args: instArgs{Arg_rd, Arg_fj}},                       // MOVFRH2GR.S rd, fj
	{mask: 0xfffffc18, value: 0x0114d800, name: "MOVGR2CF", args: instArgs{Arg_cd, Arg_rj}},                          // MOVGR2CF cd, rj
	{mask: 0xfffffc00, value: 0x0114c000, name: "MOVGR2FCSR", args: instArgs{Arg_fcsr_4_0, Arg_rj}},                  // MOVGR2FCSR fcsr, rj
	{mask: 0xfffffc00, value: 0x0114a800, name: "MOVGR2FR.D", args: instArgs{Arg_fd, Arg_rj}},                        // MOVGR2FR.D fd, rj
	{mask: 0xfffffc00, value: 0x0114a400, name: "MOVGR2FR.W", args: instArgs{Arg_fd, Arg_rj}},                        // MOVGR2FR.W fd, rj
	{mask: 0xfffffc00, value: 0x0114ac00, name: "MOVGR2FRH.W", args: instArgs{Arg_fd, Arg_rj}},                       // MOVGR2FRH.W fd, rj
	{mask: 0xffff8000, value: 0x001d8000, name: "MUL.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // MUL.D rd, rj, rk
	{mask: 0xffff8000, value: 0x001c0000, name: "MUL.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // MUL.W rd, rj, rk
	{mask: 0xffff8000, value: 0x001e0000, name: "MULH.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // MULH.D rd, rj, rk
	{mask: 0xffff8000, value: 0x001e8000, name: "MULH.DU", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                   // MULH.DU rd, rj, rk
	{mask: 0xffff8000, value: 0x001c8000, name: "MULH.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // MULH.W rd, rj, rk
	{mask: 0xffff8000, value: 0x001d0000, name: "MULH.WU", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                   // MULH.WU rd, rj, rk
	{mask: 0xffff8000, value: 0x001f0000, name: "MULW.D.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                  // MULW.D.W rd, rj, rk
	{mask: 0xffff8000, value: 0x001f8000, name: "MULW.D.WU", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                 // MULW.D.WU rd, rj, rk
	{mask: 0xffff8000, value: 0x00140000, name: "NOR", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                       // NOR rd, rj, rk
	{mask: 0xffff8000, value: 0x00150000, name: "OR", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                        // OR rd, rj, rk
	{mask: 0xffc00000, value: 0x03800000, name: "ORI", args: instArgs{Arg_rd, Arg_rj, Arg_ui12_21_10}},               // ORI rd, rj, ui12
	{mask: 0xffff8000, value: 0x00160000, name: "ORN", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                       // ORN rd, rj, rk
	{mask: 0xfe000000, value: 0x18000000, name: "PCADDI", args: instArgs{Arg_rd, Arg_si20_24_5}},                     // PCADDI rd, si20
	{mask: 0xfe000000, value: 0x1c000000, name: "PCADDU12I", args: instArgs{Arg_rd, Arg_si20_24_5}},                  // PCADDU12I rd, si20
	{mask: 0xfe000000, value: 0x1e000000, name: "PCADDU18I", args: instArgs{Arg_rd, Arg_si20_24_5}},                  // PCADDU18I rd, si20
	{mask: 0xfe000000, value: 0x1a000000, name: "PCALAU12I", args: instArgs{Arg_rd, Arg_si20_24_5}},                  // PCALAU12I rd, si20
	{mask: 0xffc00000, value: 0x2ac00000, name: "PRELD", args: instArgs{Arg_hint_4_0, Arg_rj, Arg_si12_21_10}},       // PRELD hint, rj, si12
	{mask: 0xffff8000, value: 0x382c0000, name: "PRELDX", args: instArgs{Arg_hint_4_0, Arg_rj, Arg_rk}},              // PRELDX hint, rj, rk
	{mask: 0xfffffc00, value: 0x00006800, name: "RDTIME.D", args: instArgs{Arg_rd, Arg_rj}},                          // RDTIME.D rd, rj
	{mask: 0xfffffc00, value: 0x00006400, name: "RDTIMEH.W", args: instArgs{Arg_rd, Arg_rj}},                         // RDTIMEH.W rd, rj
	{mask: 0xfffffc00, value: 0x00006000, name: "RDTIMEL.W", args: instArgs{Arg_rd, Arg_rj}},                         // RDTIMEL.W rd, rj
	{mask: 0xfffffc00, value: 0x00003000, name: "REVB.2H", args: instArgs{Arg_rd, Arg_rj}},                           // REVB.2H rd, rj
	{mask: 0xfffffc00, value: 0x00003800, name: "REVB.2W", args: instArgs{Arg_rd, Arg_rj}},                           // REVB.2W rd, rj
	{mask: 0xfffffc00, value: 0x00003400, name: "REVB.4H", args: instArgs{Arg_rd, Arg_rj}},                           // REVB.4H rd, rj
	{mask: 0xfffffc00, value: 0x00003c00, name: "REVB.D", args: instArgs{Arg_rd, Arg_rj}},                            // REVB.D rd, rj
	{mask: 0xfffffc00, value: 0x00004000, name: "REVH.2W", args: instArgs{Arg_rd, Arg_rj}},                           // REVH.2W rd, rj
	{mask: 0xfffffc00, value: 0x00004400, name: "REVH.D", args: instArgs{Arg_rd, Arg_rj}},                            // REVH.D rd, rj
	{mask: 0xffff8000, value: 0x001b8000, name: "ROTR.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // ROTR.D rd, rj, rk
	{mask: 0xffff8000, value: 0x001b0000, name: "ROTR.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // ROTR.W rd, rj, rk
	{mask: 0xffff0000, value: 0x004d0000, name: "ROTRI.D", args: instArgs{Arg_rd, Arg_rj, Arg_ui6_15_10}},            // ROTRI.D rd, rj, ui6
	{mask: 0xffff8000, value: 0x004c8000, name: "ROTRI.W", args: instArgs{Arg_rd, Arg_rj, Arg_ui5_14_10}},            // ROTRI.W rd, rj, ui5
	{mask: 0xff000000, value: 0x23000000, name: "SC.D", args: instArgs{Arg_rd, Arg_rj, Arg_si14_23_10}},              // SC.D rd, rj, si14
	{mask: 0xffff8000, value: 0x38570000, name: "SC.Q", args: instArgs{Arg_rd, Arg_rk, Arg_rj}},                      // SC.Q rd, rk, rj
	{mask: 0xff000000, value: 0x21000000, name: "SC.W", args: instArgs{Arg_rd, Arg_rj, Arg_si14_23_10}},              // SC.W rd, rj, si14
	{mask: 0xfffffc00, value: 0x38578c00, name: "SCREL.D", args: instArgs{Arg_rd, Arg_rj}},                           // SCREL.D rd, rj
	{mask: 0xfffffc00, value: 0x38578400, name: "SCREL.W", args: instArgs{Arg_rd, Arg_rj}},                           // SCREL.W rd, rj
	{mask: 0xffff8000, value: 0x00188000, name: "SLL.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // SLL.D rd, rj, rk
	{mask: 0xffff8000, value: 0x00170000, name: "SLL.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // SLL.W rd, rj, rk
	{mask: 0xffff0000, value: 0x00410000, name: "SLLI.D", args: instArgs{Arg_rd, Arg_rj, Arg_ui6_15_10}},             // SLLI.D rd, rj, ui6
	{mask: 0xffff8000, value: 0x00408000, name: "SLLI.W", args: instArgs{Arg_rd, Arg_rj, Arg_ui5_14_10}},             // SLLI.W rd, rj, ui5
	{mask: 0xffff8000, value: 0x00120000, name: "SLT", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                       // SLT rd, rj, rk
	{mask: 0xffc00000, value: 0x02000000, name: "SLTI", args: instArgs{Arg_rd, Arg_rj, Arg_si12_21_10}},              // SLTI rd, rj, si12
	{mask: 0xffff8000, value: 0x00128000, name: "SLTU", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                      // SLTU rd, rj, rk
	{mask: 0xffc00000, value: 0x02400000, name: "SLTUI", args: instArgs{Arg_rd, Arg_rj, Arg_si12_21_10}},             // SLTUI rd, rj, si12
	{mask: 0xffff8000, value: 0x00198000, name: "SRA.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // SRA.D rd, rj, rk
	{mask: 0xffff8000, value: 0x00180000, name: "SRA.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // SRA.W rd, rj, rk
	{mask: 0xffff0000, value: 0x00490000, name: "SRAI.D", args: instArgs{Arg_rd, Arg_rj, Arg_ui6_15_10}},             // SRAI.D rd, rj, ui6
	{mask: 0xffff8000, value: 0x00488000, name: "SRAI.W", args: instArgs{Arg_rd, Arg_rj, Arg_ui5_14_10}},             // SRAI.W rd, rj, ui5
	{mask: 0xffff8000, value: 0x00190000, name: "SRL.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // SRL.D rd, rj, rk
	{mask: 0xffff8000, value: 0x00178000, name: "SRL.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // SRL.W rd, rj, rk
	{mask: 0xffff0000, value: 0x00450000, name: "SRLI.D", args: instArgs{Arg_rd, Arg_rj, Arg_ui6_15_10}},             // SRLI.D rd, rj, ui6
	{mask: 0xffff8000, value: 0x00448000, name: "SRLI.W", args: instArgs{Arg_rd, Arg_rj, Arg_ui5_14_10}},             // SRLI.W rd, rj, ui5
	{mask: 0xffc00000, value: 0x29000000, name: "ST.B", args: instArgs{Arg_rd, Arg_rj, Arg_si12_21_10}},              // ST.B rd, rj, si12
	{mask: 0xffc00000, value: 0x29c00000, name: "ST.D", args: instArgs{Arg_rd, Arg_rj, Arg_si12_21_10}},              // ST.D rd, rj, si12
	{mask: 0xffc00000, value: 0x29400000, name: "ST.H", args: instArgs{Arg_rd, Arg_rj, Arg_si12_21_10}},              // ST.H rd, rj, si12
	{mask: 0xffc00000, value: 0x29800000, name: "ST.W", args: instArgs{Arg_rd, Arg_rj, Arg_si12_21_10}},              // ST.W rd, rj, si12
	{mask: 0xffff8000, value: 0x387c0000, name: "STGT.B", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // STGT.B rd, rj, rk
	{mask: 0xffff8000, value: 0x387d8000, name: "STGT.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // STGT.D rd, rj, rk
	{mask: 0xffff8000, value: 0x387c8000, name: "STGT.H", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // STGT.H rd, rj, rk
	{mask: 0xffff8000, value: 0x387d0000, name: "STGT.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // STGT.W rd, rj, rk
	{mask: 0xffff8000, value: 0x387e0000, name: "STLE.B", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // STLE.B rd, rj, rk
	{mask: 0xffff8000, value: 0x387f8000, name: "STLE.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // STLE.D rd, rj, rk
	{mask: 0xffff8000, value: 0x387e8000, name: "STLE.H", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // STLE.H rd, rj, rk
	{mask: 0xffff8000, value: 0x387f0000, name: "STLE.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                    // STLE.W rd, rj, rk
	{mask: 0xff000000, value: 0x27000000, name: "STPTR.D", args: instArgs{Arg_rd, Arg_rj, Arg_si14_23_10}},           // STPTR.D rd, rj, si14
	{mask: 0xff000000, value: 0x25000000, name: "STPTR.W", args: instArgs{Arg_rd, Arg_rj, Arg_si14_23_10}},           // STPTR.W rd, rj, si14
	{mask: 0xffff8000, value: 0x38100000, name: "STX.B", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // STX.B rd, rj, rk
	{mask: 0xffff8000, value: 0x381c0000, name: "STX.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // STX.D rd, rj, rk
	{mask: 0xffff8000, value: 0x38140000, name: "STX.H", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // STX.H rd, rj, rk
	{mask: 0xffff8000, value: 0x38180000, name: "STX.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // STX.W rd, rj, rk
	{mask: 0xffff8000, value: 0x00118000, name: "SUB.D", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // SUB.D rd, rj, rk
	{mask: 0xffff8000, value: 0x00110000, name: "SUB.W", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                     // SUB.W rd, rj, rk
	{mask: 0xffff8000, value: 0x002b0000, name: "SYSCALL", args: instArgs{Arg_code_14_0}},                            // SYSCALL code
	{mask: 0xffffffff, value: 0x06482000, name: "TLBCLR", args: instArgs{}},                                          // TLBCLR
	{mask: 0xffffffff, value: 0x06483400, name: "TLBFILL", args: instArgs{}},                                         // TLBFILL
	{mask: 0xffffffff, value: 0x06482400, name: "TLBFLUSH", args: instArgs{}},                                        // TLBFLUSH
	{mask: 0xffffffff, value: 0x06482c00, name: "TLBRD", args: instArgs{}},                                           // TLBRD
	{mask: 0xffffffff, value: 0x06482800, name: "TLBSRCH", args: instArgs{}},                                         // TLBSRCH
	{mask: 0xffffffff, value: 0x06483000, name: "TLBWR", args: instArgs{}},                                           // TLBWR
	{mask: 0xffff8000, value: 0x00158000, name: "XOR", args: instArgs{Arg_rd, Arg_rj, Arg_rk}},                       // XOR rd, rj, rk
	{mask: 0xffc00000, value: 0x03c00000, name: "XORI", args: instArgs{Arg_rd, Arg_rj, Arg_ui12_21_10}},              // XORI rd, rj, ui12
}

// 指令编码信息表
var AOpContextTable_comment = [...]string{
	"// ADD.D rd, rj, rk",
	"// ADD.W rd, rj, rk",
	"// ADDI.D rd, rj, si12",
	"// ADDI.W rd, rj, si12",
	"// ADDU16I.D rd, rj, si16",
	"// ALSL.D rd, rj, rk, sa2",
	"// ALSL.W rd, rj, rk, sa2",
	"// ALSL.WU rd, rj, rk, sa2",
	"// AMADD.B rd, rk, rj",
	"// AMADD.D rd, rk, rj",
	"// AMADD.H rd, rk, rj",
	"// AMADD.W rd, rk, rj",
	"// AMADD_DB.B rd, rk, rj",
	"// AMADD_DB.D rd, rk, rj",
	"// AMADD_DB.H rd, rk, rj",
	"// AMADD_DB.W rd, rk, rj",
	"// AMAND.D rd, rk, rj",
	"// AMAND.W rd, rk, rj",
	"// AMAND_DB.D rd, rk, rj",
	"// AMAND_DB.W rd, rk, rj",
	"// AMCAS.B rd, rk, rj",
	"// AMCAS.D rd, rk, rj",
	"// AMCAS.H rd, rk, rj",
	"// AMCAS.W rd, rk, rj",
	"// AMCAS_DB.B rd, rk, rj",
	"// AMCAS_DB.D rd, rk, rj",
	"// AMCAS_DB.H rd, rk, rj",
	"// AMCAS_DB.W rd, rk, rj",
	"// AMMAX.D rd, rk, rj",
	"// AMMAX.DU rd, rk, rj",
	"// AMMAX.W rd, rk, rj",
	"// AMMAX.WU rd, rk, rj",
	"// AMMAX_DB.D rd, rk, rj",
	"// AMMAX_DB.DU rd, rk, rj",
	"// AMMAX_DB.W rd, rk, rj",
	"// AMMAX_DB.WU rd, rk, rj",
	"// AMMIN.D rd, rk, rj",
	"// AMMIN.DU rd, rk, rj",
	"// AMMIN.W rd, rk, rj",
	"// AMMIN.WU rd, rk, rj",
	"// AMMIN_DB.D rd, rk, rj",
	"// AMMIN_DB.DU rd, rk, rj",
	"// AMMIN_DB.W rd, rk, rj",
	"// AMMIN_DB.WU rd, rk, rj",
	"// AMOR.D rd, rk, rj",
	"// AMOR.W rd, rk, rj",
	"// AMOR_DB.D rd, rk, rj",
	"// AMOR_DB.W rd, rk, rj",
	"// AMSWAP.B rd, rk, rj",
	"// AMSWAP.D rd, rk, rj",
	"// AMSWAP.H rd, rk, rj",
	"// AMSWAP.W rd, rk, rj",
	"// AMSWAP_DB.B rd, rk, rj",
	"// AMSWAP_DB.D rd, rk, rj",
	"// AMSWAP_DB.H rd, rk, rj",
	"// AMSWAP_DB.W rd, rk, rj",
	"// AMXOR.D rd, rk, rj",
	"// AMXOR.W rd, rk, rj",
	"// AMXOR_DB.D rd, rk, rj",
	"// AMXOR_DB.W rd, rk, rj",
	"// AND rd, rj, rk",
	"// ANDI rd, rj, ui12",
	"// ANDN rd, rj, rk",
	"// ASRTGT.D rj, rk",
	"// ASRTLE.D rj, rk",
	"// B offs",
	"// BCEQZ cj, offs",
	"// BCNEZ cj, offs",
	"// BEQ rj, rd, offs",
	"// BEQZ rj, offs",
	"// BGE rj, rd, offs",
	"// BGEU rj, rd, offs",
	"// BITREV.4B rd, rj",
	"// BITREV.8B rd, rj",
	"// BITREV.D rd, rj",
	"// BITREV.W rd, rj",
	"// BL offs",
	"// BLT rj, rd, offs",
	"// BLTU rj, rd, offs",
	"// BNE rj, rd, offs",
	"// BNEZ rj, offs",
	"// BREAK code",
	"// BSTRINS.D rd, rj, msbd, lsbd",
	"// BSTRINS.W rd, rj, msbw, lsbw",
	"// BSTRPICK.D rd, rj, msbd, lsbd",
	"// BSTRPICK.W rd, rj, msbw, lsbw",
	"// BYTEPICK.D rd, rj, rk, sa3",
	"// BYTEPICK.W rd, rj, rk, sa2",
	"// CACOP code, rj, si12",
	"// CLO.D rd, rj",
	"// CLO.W rd, rj",
	"// CLZ.D rd, rj",
	"// CLZ.W rd, rj",
	"// CPUCFG rd, rj",
	"// CRC.W.B.W rd, rj, rk",
	"// CRC.W.D.W rd, rj, rk",
	"// CRC.W.H.W rd, rj, rk",
	"// CRC.W.W.W rd, rj, rk",
	"// CRCC.W.B.W rd, rj, rk",
	"// CRCC.W.D.W rd, rj, rk",
	"// CRCC.W.H.W rd, rj, rk",
	"// CRCC.W.W.W rd, rj, rk",
	"// CSRRD rd, csr",
	"// CSRWR rd, csr",
	"// CSRXCHG rd, rj, csr",
	"// CTO.D rd, rj",
	"// CTO.W rd, rj",
	"// CTZ.D rd, rj",
	"// CTZ.W rd, rj",
	"// DBAR hint",
	"// DBCL code",
	"// DIV.D rd, rj, rk",
	"// DIV.DU rd, rj, rk",
	"// DIV.W rd, rj, rk",
	"// DIV.WU rd, rj, rk",
	"// ERTN",
	"// EXT.W.B rd, rj",
	"// EXT.W.H rd, rj",
	"// FABS.D fd, fj",
	"// FABS.S fd, fj",
	"// FADD.D fd, fj, fk",
	"// FADD.S fd, fj, fk",
	"// FCLASS.D fd, fj",
	"// FCLASS.S fd, fj",
	"// FCMP.CAF.D cd, fj, fk",
	"// FCMP.CAF.S cd, fj, fk",
	"// FCMP.CEQ.D cd, fj, fk",
	"// FCMP.CEQ.S cd, fj, fk",
	"// FCMP.CLE.D cd, fj, fk",
	"// FCMP.CLE.S cd, fj, fk",
	"// FCMP.CLT.D cd, fj, fk",
	"// FCMP.CLT.S cd, fj, fk",
	"// FCMP.CNE.D cd, fj, fk",
	"// FCMP.CNE.S cd, fj, fk",
	"// FCMP.COR.D cd, fj, fk",
	"// FCMP.COR.S cd, fj, fk",
	"// FCMP.CUEQ.D cd, fj, fk",
	"// FCMP.CUEQ.S cd, fj, fk",
	"// FCMP.CULE.D cd, fj, fk",
	"// FCMP.CULE.S cd, fj, fk",
	"// FCMP.CULT.D cd, fj, fk",
	"// FCMP.CULT.S cd, fj, fk",
	"// FCMP.CUN.D cd, fj, fk",
	"// FCMP.CUN.S cd, fj, fk",
	"// FCMP.CUNE.D cd, fj, fk",
	"// FCMP.CUNE.S cd, fj, fk",
	"// FCMP.SAF.D cd, fj, fk",
	"// FCMP.SAF.S cd, fj, fk",
	"// FCMP.SEQ.D cd, fj, fk",
	"// FCMP.SEQ.S cd, fj, fk",
	"// FCMP.SLE.D cd, fj, fk",
	"// FCMP.SLE.S cd, fj, fk",
	"// FCMP.SLT.D cd, fj, fk",
	"// FCMP.SLT.S cd, fj, fk",
	"// FCMP.SNE.D cd, fj, fk",
	"// FCMP.SNE.S cd, fj, fk",
	"// FCMP.SOR.D cd, fj, fk",
	"// FCMP.SOR.S cd, fj, fk",
	"// FCMP.SUEQ.D cd, fj, fk",
	"// FCMP.SUEQ.S cd, fj, fk",
	"// FCMP.SULE.D cd, fj, fk",
	"// FCMP.SULE.S cd, fj, fk",
	"// FCMP.SULT.D cd, fj, fk",
	"// FCMP.SULT.S cd, fj, fk",
	"// FCMP.SUN.D cd, fj, fk",
	"// FCMP.SUN.S cd, fj, fk",
	"// FCMP.SUNE.D cd, fj, fk",
	"// FCMP.SUNE.S cd, fj, fk",
	"// FCOPYSIGN.D fd, fj, fk",
	"// FCOPYSIGN.S fd, fj, fk",
	"// FCVT.D.S fd, fj",
	"// FCVT.S.D fd, fj",
	"// FDIV.D fd, fj, fk",
	"// FDIV.S fd, fj, fk",
	"// FFINT.D.L fd, fj",
	"// FFINT.D.W fd, fj",
	"// FFINT.S.L fd, fj",
	"// FFINT.S.W fd, fj",
	"// FLD.D fd, rj, si12",
	"// FLD.S fd, rj, si12",
	"// FLDGT.D fd, rj, rk",
	"// FLDGT.S fd, rj, rk",
	"// FLDLE.D fd, rj, rk",
	"// FLDLE.S fd, rj, rk",
	"// FLDX.D fd, rj, rk",
	"// FLDX.S fd, rj, rk",
	"// FLOGB.D fd, fj",
	"// FLOGB.S fd, fj",
	"// FMADD.D fd, fj, fk, fa",
	"// FMADD.S fd, fj, fk, fa",
	"// FMAX.D fd, fj, fk",
	"// FMAX.S fd, fj, fk",
	"// FMAXA.D fd, fj, fk",
	"// FMAXA.S fd, fj, fk",
	"// FMIN.D fd, fj, fk",
	"// FMIN.S fd, fj, fk",
	"// FMINA.D fd, fj, fk",
	"// FMINA.S fd, fj, fk",
	"// FMOV.D fd, fj",
	"// FMOV.S fd, fj",
	"// FMSUB.D fd, fj, fk, fa",
	"// FMSUB.S fd, fj, fk, fa",
	"// FMUL.D fd, fj, fk",
	"// FMUL.S fd, fj, fk",
	"// FNEG.D fd, fj",
	"// FNEG.S fd, fj",
	"// FNMADD.D fd, fj, fk, fa",
	"// FNMADD.S fd, fj, fk, fa",
	"// FNMSUB.D fd, fj, fk, fa",
	"// FNMSUB.S fd, fj, fk, fa",
	"// FRECIP.D fd, fj",
	"// FRECIP.S fd, fj",
	"// FRECIPE.D fd, fj",
	"// FRECIPE.S fd, fj",
	"// FRINT.D fd, fj",
	"// FRINT.S fd, fj",
	"// FRSQRT.D fd, fj",
	"// FRSQRT.S fd, fj",
	"// FRSQRTE.D fd, fj",
	"// FRSQRTE.S fd, fj",
	"// FSCALEB.D fd, fj, fk",
	"// FSCALEB.S fd, fj, fk",
	"// FSEL fd, fj, fk, ca",
	"// FSQRT.D fd, fj",
	"// FSQRT.S fd, fj",
	"// FST.D fd, rj, si12",
	"// FST.S fd, rj, si12",
	"// FSTGT.D fd, rj, rk",
	"// FSTGT.S fd, rj, rk",
	"// FSTLE.D fd, rj, rk",
	"// FSTLE.S fd, rj, rk",
	"// FSTX.D fd, rj, rk",
	"// FSTX.S fd, rj, rk",
	"// FSUB.D fd, fj, fk",
	"// FSUB.S fd, fj, fk",
	"// FTINT.L.D fd, fj",
	"// FTINT.L.S fd, fj",
	"// FTINT.W.D fd, fj",
	"// FTINT.W.S fd, fj",
	"// FTINTRM.L.D fd, fj",
	"// FTINTRM.L.S fd, fj",
	"// FTINTRM.W.D fd, fj",
	"// FTINTRM.W.S fd, fj",
	"// FTINTRNE.L.D fd, fj",
	"// FTINTRNE.L.S fd, fj",
	"// FTINTRNE.W.D fd, fj",
	"// FTINTRNE.W.S fd, fj",
	"// FTINTRP.L.D fd, fj",
	"// FTINTRP.L.S fd, fj",
	"// FTINTRP.W.D fd, fj",
	"// FTINTRP.W.S fd, fj",
	"// FTINTRZ.L.D fd, fj",
	"// FTINTRZ.L.S fd, fj",
	"// FTINTRZ.W.D fd, fj",
	"// FTINTRZ.W.S fd, fj",
	"// IBAR hint",
	"// IDLE level",
	"// INVTLB op, rj, rk",
	"// IOCSRRD.B rd, rj",
	"// IOCSRRD.D rd, rj",
	"// IOCSRRD.H rd, rj",
	"// IOCSRRD.W rd, rj",
	"// IOCSRWR.B rd, rj",
	"// IOCSRWR.D rd, rj",
	"// IOCSRWR.H rd, rj",
	"// IOCSRWR.W rd, rj",
	"// JIRL rd, rj, offs",
	"// LD.B rd, rj, si12",
	"// LD.BU rd, rj, si12",
	"// LD.D rd, rj, si12",
	"// LD.H rd, rj, si12",
	"// LD.HU rd, rj, si12",
	"// LD.W rd, rj, si12",
	"// LD.WU rd, rj, si12",
	"// LDDIR rd, rj, level",
	"// LDGT.B rd, rj, rk",
	"// LDGT.D rd, rj, rk",
	"// LDGT.H rd, rj, rk",
	"// LDGT.W rd, rj, rk",
	"// LDLE.B rd, rj, rk",
	"// LDLE.D rd, rj, rk",
	"// LDLE.H rd, rj, rk",
	"// LDLE.W rd, rj, rk",
	"// LDPTE rj, seq",
	"// LDPTR.D rd, rj, si14",
	"// LDPTR.W rd, rj, si14",
	"// LDX.B rd, rj, rk",
	"// LDX.BU rd, rj, rk",
	"// LDX.D rd, rj, rk",
	"// LDX.H rd, rj, rk",
	"// LDX.HU rd, rj, rk",
	"// LDX.W rd, rj, rk",
	"// LDX.WU rd, rj, rk",
	"// LL.D rd, rj, si14",
	"// LL.W rd, rj, si14",
	"// LLACQ.D rd, rj",
	"// LLACQ.W rd, rj",
	"// LU12I.W rd, si20",
	"// LU32I.D rd, si20",
	"// LU52I.D rd, rj, si12",
	"// MASKEQZ rd, rj, rk",
	"// MASKNEZ rd, rj, rk",
	"// MOD.D rd, rj, rk",
	"// MOD.DU rd, rj, rk",
	"// MOD.W rd, rj, rk",
	"// MOD.WU rd, rj, rk",
	"// MOVCF2FR fd, cj",
	"// MOVCF2GR rd, cj",
	"// MOVFCSR2GR rd, fcsr",
	"// MOVFR2CF cd, fj",
	"// MOVFR2GR.D rd, fj",
	"// MOVFR2GR.S rd, fj",
	"// MOVFRH2GR.S rd, fj",
	"// MOVGR2CF cd, rj",
	"// MOVGR2FCSR fcsr, rj",
	"// MOVGR2FR.D fd, rj",
	"// MOVGR2FR.W fd, rj",
	"// MOVGR2FRH.W fd, rj",
	"// MUL.D rd, rj, rk",
	"// MUL.W rd, rj, rk",
	"// MULH.D rd, rj, rk",
	"// MULH.DU rd, rj, rk",
	"// MULH.W rd, rj, rk",
	"// MULH.WU rd, rj, rk",
	"// MULW.D.W rd, rj, rk",
	"// MULW.D.WU rd, rj, rk",
	"// NOR rd, rj, rk",
	"// OR rd, rj, rk",
	"// ORI rd, rj, ui12",
	"// ORN rd, rj, rk",
	"// PCADDI rd, si20",
	"// PCADDU12I rd, si20",
	"// PCADDU18I rd, si20",
	"// PCALAU12I rd, si20",
	"// PRELD hint, rj, si12",
	"// PRELDX hint, rj, rk",
	"// RDTIME.D rd, rj",
	"// RDTIMEH.W rd, rj",
	"// RDTIMEL.W rd, rj",
	"// REVB.2H rd, rj",
	"// REVB.2W rd, rj",
	"// REVB.4H rd, rj",
	"// REVB.D rd, rj",
	"// REVH.2W rd, rj",
	"// REVH.D rd, rj",
	"// ROTR.D rd, rj, rk",
	"// ROTR.W rd, rj, rk",
	"// ROTRI.D rd, rj, ui6",
	"// ROTRI.W rd, rj, ui5",
	"// SC.D rd, rj, si14",
	"// SC.Q rd, rk, rj",
	"// SC.W rd, rj, si14",
	"// SCREL.D rd, rj",
	"// SCREL.W rd, rj",
	"// SLL.D rd, rj, rk",
	"// SLL.W rd, rj, rk",
	"// SLLI.D rd, rj, ui6",
	"// SLLI.W rd, rj, ui5",
	"// SLT rd, rj, rk",
	"// SLTI rd, rj, si12",
	"// SLTU rd, rj, rk",
	"// SLTUI rd, rj, si12",
	"// SRA.D rd, rj, rk",
	"// SRA.W rd, rj, rk",
	"// SRAI.D rd, rj, ui6",
	"// SRAI.W rd, rj, ui5",
	"// SRL.D rd, rj, rk",
	"// SRL.W rd, rj, rk",
	"// SRLI.D rd, rj, ui6",
	"// SRLI.W rd, rj, ui5",
	"// ST.B rd, rj, si12",
	"// ST.D rd, rj, si12",
	"// ST.H rd, rj, si12",
	"// ST.W rd, rj, si12",
	"// STGT.B rd, rj, rk",
	"// STGT.D rd, rj, rk",
	"// STGT.H rd, rj, rk",
	"// STGT.W rd, rj, rk",
	"// STLE.B rd, rj, rk",
	"// STLE.D rd, rj, rk",
	"// STLE.H rd, rj, rk",
	"// STLE.W rd, rj, rk",
	"// STPTR.D rd, rj, si14",
	"// STPTR.W rd, rj, si14",
	"// STX.B rd, rj, rk",
	"// STX.D rd, rj, rk",
	"// STX.H rd, rj, rk",
	"// STX.W rd, rj, rk",
	"// SUB.D rd, rj, rk",
	"// SUB.W rd, rj, rk",
	"// SYSCALL code",
	"// TLBCLR",
	"// TLBFILL",
	"// TLBFLUSH",
	"// TLBRD",
	"// TLBSRCH",
	"// TLBWR",
	"// XOR rd, rj, rk",
	"// XORI rd, rj, ui12",
}
