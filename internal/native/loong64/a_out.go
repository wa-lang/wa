// Copyright (C) 2025 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

// 注意: 此代码是程序生成, 不要手动修改!!!

package loong64

import "wa-lang.org/wa/internal/native/abi"

const (
	_             abi.As = iota
	AADDI_D              // ADDI.D rd, rj, si12
	AADDI_W              // ADDI.W rd, rj, si12
	AADDU16I_D           // ADDU16I.D rd, rj, si16
	AADD_D               // ADD.D rd, rj, rk
	AADD_W               // ADD.W rd, rj, rk
	AALSL_D              // ALSL.D rd, rj, rk, sa2
	AALSL_W              // ALSL.W rd, rj, rk, sa2
	AALSL_WU             // ALSL.WU rd, rj, rk, sa2
	AAMADD_B             // AMADD.B rd, rk, rj
	AAMADD_D             // AMADD.D rd, rk, rj
	AAMADD_DB_B          // AMADD_DB.B rd, rk, rj
	AAMADD_DB_D          // AMADD_DB.D rd, rk, rj
	AAMADD_DB_H          // AMADD_DB.H rd, rk, rj
	AAMADD_DB_W          // AMADD_DB.W rd, rk, rj
	AAMADD_H             // AMADD.H rd, rk, rj
	AAMADD_W             // AMADD.W rd, rk, rj
	AAMAND_D             // AMAND.D rd, rk, rj
	AAMAND_DB_D          // AMAND_DB.D rd, rk, rj
	AAMAND_DB_W          // AMAND_DB.W rd, rk, rj
	AAMAND_W             // AMAND.W rd, rk, rj
	AAMCAS_B             // AMCAS.B rd, rk, rj
	AAMCAS_D             // AMCAS.D rd, rk, rj
	AAMCAS_DB_B          // AMCAS_DB.B rd, rk, rj
	AAMCAS_DB_D          // AMCAS_DB.D rd, rk, rj
	AAMCAS_DB_H          // AMCAS_DB.H rd, rk, rj
	AAMCAS_DB_W          // AMCAS_DB.W rd, rk, rj
	AAMCAS_H             // AMCAS.H rd, rk, rj
	AAMCAS_W             // AMCAS.W rd, rk, rj
	AAMMAX_D             // AMMAX.D rd, rk, rj
	AAMMAX_DB_D          // AMMAX_DB.D rd, rk, rj
	AAMMAX_DB_DU         // AMMAX_DB.DU rd, rk, rj
	AAMMAX_DB_W          // AMMAX_DB.W rd, rk, rj
	AAMMAX_DB_WU         // AMMAX_DB.WU rd, rk, rj
	AAMMAX_DU            // AMMAX.DU rd, rk, rj
	AAMMAX_W             // AMMAX.W rd, rk, rj
	AAMMAX_WU            // AMMAX.WU rd, rk, rj
	AAMMIN_D             // AMMIN.D rd, rk, rj
	AAMMIN_DB_D          // AMMIN_DB.D rd, rk, rj
	AAMMIN_DB_DU         // AMMIN_DB.DU rd, rk, rj
	AAMMIN_DB_W          // AMMIN_DB.W rd, rk, rj
	AAMMIN_DB_WU         // AMMIN_DB.WU rd, rk, rj
	AAMMIN_DU            // AMMIN.DU rd, rk, rj
	AAMMIN_W             // AMMIN.W rd, rk, rj
	AAMMIN_WU            // AMMIN.WU rd, rk, rj
	AAMOR_D              // AMOR.D rd, rk, rj
	AAMOR_DB_D           // AMOR_DB.D rd, rk, rj
	AAMOR_DB_W           // AMOR_DB.W rd, rk, rj
	AAMOR_W              // AMOR.W rd, rk, rj
	AAMSWAP_B            // AMSWAP.B rd, rk, rj
	AAMSWAP_D            // AMSWAP.D rd, rk, rj
	AAMSWAP_DB_B         // AMSWAP_DB.B rd, rk, rj
	AAMSWAP_DB_D         // AMSWAP_DB.D rd, rk, rj
	AAMSWAP_DB_H         // AMSWAP_DB.H rd, rk, rj
	AAMSWAP_DB_W         // AMSWAP_DB.W rd, rk, rj
	AAMSWAP_H            // AMSWAP.H rd, rk, rj
	AAMSWAP_W            // AMSWAP.W rd, rk, rj
	AAMXOR_D             // AMXOR.D rd, rk, rj
	AAMXOR_DB_D          // AMXOR_DB.D rd, rk, rj
	AAMXOR_DB_W          // AMXOR_DB.W rd, rk, rj
	AAMXOR_W             // AMXOR.W rd, rk, rj
	AAND                 // AND rd, rj, rk
	AANDI                // ANDI rd, rj, ui12
	AANDN                // ANDN rd, rj, rk
	AASRTGT_D            // ASRTGT.D rj, rk
	AASRTLE_D            // ASRTLE.D rj, rk
	AB                   // B offs
	ABCEQZ               // BCEQZ cj, offs
	ABCNEZ               // BCNEZ cj, offs
	ABEQ                 // BEQ rj, rd, offs
	ABEQZ                // BEQZ rj, offs
	ABGE                 // BGE rj, rd, offs
	ABGEU                // BGEU rj, rd, offs
	ABITREV_4B           // BITREV.4B rd, rj
	ABITREV_8B           // BITREV.8B rd, rj
	ABITREV_D            // BITREV.D rd, rj
	ABITREV_W            // BITREV.W rd, rj
	ABL                  // BL offs
	ABLT                 // BLT rj, rd, offs
	ABLTU                // BLTU rj, rd, offs
	ABNE                 // BNE rj, rd, offs
	ABNEZ                // BNEZ rj, offs
	ABREAK               // BREAK code
	ABSTRINS_D           // BSTRINS.D rd, rj, msbd, lsbd
	ABSTRINS_W           // BSTRINS.W rd, rj, msbw, lsbw
	ABSTRPICK_D          // BSTRPICK.D rd, rj, msbd, lsbd
	ABSTRPICK_W          // BSTRPICK.W rd, rj, msbw, lsbw
	ABYTEPICK_D          // BYTEPICK.D rd, rj, rk, sa3
	ABYTEPICK_W          // BYTEPICK.W rd, rj, rk, sa2
	ACACOP               // CACOP code, rj, si12
	ACLO_D               // CLO.D rd, rj
	ACLO_W               // CLO.W rd, rj
	ACLZ_D               // CLZ.D rd, rj
	ACLZ_W               // CLZ.W rd, rj
	ACPUCFG              // CPUCFG rd, rj
	ACRCC_W_B_W          // CRCC.W.B.W rd, rj, rk
	ACRCC_W_D_W          // CRCC.W.D.W rd, rj, rk
	ACRCC_W_H_W          // CRCC.W.H.W rd, rj, rk
	ACRCC_W_W_W          // CRCC.W.W.W rd, rj, rk
	ACRC_W_B_W           // CRC.W.B.W rd, rj, rk
	ACRC_W_D_W           // CRC.W.D.W rd, rj, rk
	ACRC_W_H_W           // CRC.W.H.W rd, rj, rk
	ACRC_W_W_W           // CRC.W.W.W rd, rj, rk
	ACSRRD               // CSRRD rd, csr
	ACSRWR               // CSRWR rd, csr
	ACSRXCHG             // CSRXCHG rd, rj, csr
	ACTO_D               // CTO.D rd, rj
	ACTO_W               // CTO.W rd, rj
	ACTZ_D               // CTZ.D rd, rj
	ACTZ_W               // CTZ.W rd, rj
	ADBAR                // DBAR hint
	ADBCL                // DBCL code
	ADIV_D               // DIV.D rd, rj, rk
	ADIV_DU              // DIV.DU rd, rj, rk
	ADIV_W               // DIV.W rd, rj, rk
	ADIV_WU              // DIV.WU rd, rj, rk
	AERTN                // ERTN
	AEXT_W_B             // EXT.W.B rd, rj
	AEXT_W_H             // EXT.W.H rd, rj
	AFABS_D              // FABS.D fd, fj
	AFABS_S              // FABS.S fd, fj
	AFADD_D              // FADD.D fd, fj, fk
	AFADD_S              // FADD.S fd, fj, fk
	AFCLASS_D            // FCLASS.D fd, fj
	AFCLASS_S            // FCLASS.S fd, fj
	AFCMP_CAF_D          // FCMP.CAF.D cd, fj, fk
	AFCMP_CAF_S          // FCMP.CAF.S cd, fj, fk
	AFCMP_CEQ_D          // FCMP.CEQ.D cd, fj, fk
	AFCMP_CEQ_S          // FCMP.CEQ.S cd, fj, fk
	AFCMP_CLE_D          // FCMP.CLE.D cd, fj, fk
	AFCMP_CLE_S          // FCMP.CLE.S cd, fj, fk
	AFCMP_CLT_D          // FCMP.CLT.D cd, fj, fk
	AFCMP_CLT_S          // FCMP.CLT.S cd, fj, fk
	AFCMP_CNE_D          // FCMP.CNE.D cd, fj, fk
	AFCMP_CNE_S          // FCMP.CNE.S cd, fj, fk
	AFCMP_COR_D          // FCMP.COR.D cd, fj, fk
	AFCMP_COR_S          // FCMP.COR.S cd, fj, fk
	AFCMP_CUEQ_D         // FCMP.CUEQ.D cd, fj, fk
	AFCMP_CUEQ_S         // FCMP.CUEQ.S cd, fj, fk
	AFCMP_CULE_D         // FCMP.CULE.D cd, fj, fk
	AFCMP_CULE_S         // FCMP.CULE.S cd, fj, fk
	AFCMP_CULT_D         // FCMP.CULT.D cd, fj, fk
	AFCMP_CULT_S         // FCMP.CULT.S cd, fj, fk
	AFCMP_CUNE_D         // FCMP.CUNE.D cd, fj, fk
	AFCMP_CUNE_S         // FCMP.CUNE.S cd, fj, fk
	AFCMP_CUN_D          // FCMP.CUN.D cd, fj, fk
	AFCMP_CUN_S          // FCMP.CUN.S cd, fj, fk
	AFCMP_SAF_D          // FCMP.SAF.D cd, fj, fk
	AFCMP_SAF_S          // FCMP.SAF.S cd, fj, fk
	AFCMP_SEQ_D          // FCMP.SEQ.D cd, fj, fk
	AFCMP_SEQ_S          // FCMP.SEQ.S cd, fj, fk
	AFCMP_SLE_D          // FCMP.SLE.D cd, fj, fk
	AFCMP_SLE_S          // FCMP.SLE.S cd, fj, fk
	AFCMP_SLT_D          // FCMP.SLT.D cd, fj, fk
	AFCMP_SLT_S          // FCMP.SLT.S cd, fj, fk
	AFCMP_SNE_D          // FCMP.SNE.D cd, fj, fk
	AFCMP_SNE_S          // FCMP.SNE.S cd, fj, fk
	AFCMP_SOR_D          // FCMP.SOR.D cd, fj, fk
	AFCMP_SOR_S          // FCMP.SOR.S cd, fj, fk
	AFCMP_SUEQ_D         // FCMP.SUEQ.D cd, fj, fk
	AFCMP_SUEQ_S         // FCMP.SUEQ.S cd, fj, fk
	AFCMP_SULE_D         // FCMP.SULE.D cd, fj, fk
	AFCMP_SULE_S         // FCMP.SULE.S cd, fj, fk
	AFCMP_SULT_D         // FCMP.SULT.D cd, fj, fk
	AFCMP_SULT_S         // FCMP.SULT.S cd, fj, fk
	AFCMP_SUNE_D         // FCMP.SUNE.D cd, fj, fk
	AFCMP_SUNE_S         // FCMP.SUNE.S cd, fj, fk
	AFCMP_SUN_D          // FCMP.SUN.D cd, fj, fk
	AFCMP_SUN_S          // FCMP.SUN.S cd, fj, fk
	AFCOPYSIGN_D         // FCOPYSIGN.D fd, fj, fk
	AFCOPYSIGN_S         // FCOPYSIGN.S fd, fj, fk
	AFCVT_D_S            // FCVT.D.S fd, fj
	AFCVT_S_D            // FCVT.S.D fd, fj
	AFDIV_D              // FDIV.D fd, fj, fk
	AFDIV_S              // FDIV.S fd, fj, fk
	AFFINT_D_L           // FFINT.D.L fd, fj
	AFFINT_D_W           // FFINT.D.W fd, fj
	AFFINT_S_L           // FFINT.S.L fd, fj
	AFFINT_S_W           // FFINT.S.W fd, fj
	AFLDGT_D             // FLDGT.D fd, rj, rk
	AFLDGT_S             // FLDGT.S fd, rj, rk
	AFLDLE_D             // FLDLE.D fd, rj, rk
	AFLDLE_S             // FLDLE.S fd, rj, rk
	AFLDX_D              // FLDX.D fd, rj, rk
	AFLDX_S              // FLDX.S fd, rj, rk
	AFLD_D               // FLD.D fd, rj, si12
	AFLD_S               // FLD.S fd, rj, si12
	AFLOGB_D             // FLOGB.D fd, fj
	AFLOGB_S             // FLOGB.S fd, fj
	AFMADD_D             // FMADD.D fd, fj, fk, fa
	AFMADD_S             // FMADD.S fd, fj, fk, fa
	AFMAXA_D             // FMAXA.D fd, fj, fk
	AFMAXA_S             // FMAXA.S fd, fj, fk
	AFMAX_D              // FMAX.D fd, fj, fk
	AFMAX_S              // FMAX.S fd, fj, fk
	AFMINA_D             // FMINA.D fd, fj, fk
	AFMINA_S             // FMINA.S fd, fj, fk
	AFMIN_D              // FMIN.D fd, fj, fk
	AFMIN_S              // FMIN.S fd, fj, fk
	AFMOV_D              // FMOV.D fd, fj
	AFMOV_S              // FMOV.S fd, fj
	AFMSUB_D             // FMSUB.D fd, fj, fk, fa
	AFMSUB_S             // FMSUB.S fd, fj, fk, fa
	AFMUL_D              // FMUL.D fd, fj, fk
	AFMUL_S              // FMUL.S fd, fj, fk
	AFNEG_D              // FNEG.D fd, fj
	AFNEG_S              // FNEG.S fd, fj
	AFNMADD_D            // FNMADD.D fd, fj, fk, fa
	AFNMADD_S            // FNMADD.S fd, fj, fk, fa
	AFNMSUB_D            // FNMSUB.D fd, fj, fk, fa
	AFNMSUB_S            // FNMSUB.S fd, fj, fk, fa
	AFRECIPE_D           // FRECIPE.D fd, fj
	AFRECIPE_S           // FRECIPE.S fd, fj
	AFRECIP_D            // FRECIP.D fd, fj
	AFRECIP_S            // FRECIP.S fd, fj
	AFRINT_D             // FRINT.D fd, fj
	AFRINT_S             // FRINT.S fd, fj
	AFRSQRTE_D           // FRSQRTE.D fd, fj
	AFRSQRTE_S           // FRSQRTE.S fd, fj
	AFRSQRT_D            // FRSQRT.D fd, fj
	AFRSQRT_S            // FRSQRT.S fd, fj
	AFSCALEB_D           // FSCALEB.D fd, fj, fk
	AFSCALEB_S           // FSCALEB.S fd, fj, fk
	AFSEL                // FSEL fd, fj, fk, ca
	AFSQRT_D             // FSQRT.D fd, fj
	AFSQRT_S             // FSQRT.S fd, fj
	AFSTGT_D             // FSTGT.D fd, rj, rk
	AFSTGT_S             // FSTGT.S fd, rj, rk
	AFSTLE_D             // FSTLE.D fd, rj, rk
	AFSTLE_S             // FSTLE.S fd, rj, rk
	AFSTX_D              // FSTX.D fd, rj, rk
	AFSTX_S              // FSTX.S fd, rj, rk
	AFST_D               // FST.D fd, rj, si12
	AFST_S               // FST.S fd, rj, si12
	AFSUB_D              // FSUB.D fd, fj, fk
	AFSUB_S              // FSUB.S fd, fj, fk
	AFTINTRM_L_D         // FTINTRM.L.D fd, fj
	AFTINTRM_L_S         // FTINTRM.L.S fd, fj
	AFTINTRM_W_D         // FTINTRM.W.D fd, fj
	AFTINTRM_W_S         // FTINTRM.W.S fd, fj
	AFTINTRNE_L_D        // FTINTRNE.L.D fd, fj
	AFTINTRNE_L_S        // FTINTRNE.L.S fd, fj
	AFTINTRNE_W_D        // FTINTRNE.W.D fd, fj
	AFTINTRNE_W_S        // FTINTRNE.W.S fd, fj
	AFTINTRP_L_D         // FTINTRP.L.D fd, fj
	AFTINTRP_L_S         // FTINTRP.L.S fd, fj
	AFTINTRP_W_D         // FTINTRP.W.D fd, fj
	AFTINTRP_W_S         // FTINTRP.W.S fd, fj
	AFTINTRZ_L_D         // FTINTRZ.L.D fd, fj
	AFTINTRZ_L_S         // FTINTRZ.L.S fd, fj
	AFTINTRZ_W_D         // FTINTRZ.W.D fd, fj
	AFTINTRZ_W_S         // FTINTRZ.W.S fd, fj
	AFTINT_L_D           // FTINT.L.D fd, fj
	AFTINT_L_S           // FTINT.L.S fd, fj
	AFTINT_W_D           // FTINT.W.D fd, fj
	AFTINT_W_S           // FTINT.W.S fd, fj
	AIBAR                // IBAR hint
	AIDLE                // IDLE level
	AINVTLB              // INVTLB op, rj, rk
	AIOCSRRD_B           // IOCSRRD.B rd, rj
	AIOCSRRD_D           // IOCSRRD.D rd, rj
	AIOCSRRD_H           // IOCSRRD.H rd, rj
	AIOCSRRD_W           // IOCSRRD.W rd, rj
	AIOCSRWR_B           // IOCSRWR.B rd, rj
	AIOCSRWR_D           // IOCSRWR.D rd, rj
	AIOCSRWR_H           // IOCSRWR.H rd, rj
	AIOCSRWR_W           // IOCSRWR.W rd, rj
	AJIRL                // JIRL rd, rj, offs
	ALDDIR               // LDDIR rd, rj, level
	ALDGT_B              // LDGT.B rd, rj, rk
	ALDGT_D              // LDGT.D rd, rj, rk
	ALDGT_H              // LDGT.H rd, rj, rk
	ALDGT_W              // LDGT.W rd, rj, rk
	ALDLE_B              // LDLE.B rd, rj, rk
	ALDLE_D              // LDLE.D rd, rj, rk
	ALDLE_H              // LDLE.H rd, rj, rk
	ALDLE_W              // LDLE.W rd, rj, rk
	ALDPTE               // LDPTE rj, seq
	ALDPTR_D             // LDPTR.D rd, rj, si14
	ALDPTR_W             // LDPTR.W rd, rj, si14
	ALDX_B               // LDX.B rd, rj, rk
	ALDX_BU              // LDX.BU rd, rj, rk
	ALDX_D               // LDX.D rd, rj, rk
	ALDX_H               // LDX.H rd, rj, rk
	ALDX_HU              // LDX.HU rd, rj, rk
	ALDX_W               // LDX.W rd, rj, rk
	ALDX_WU              // LDX.WU rd, rj, rk
	ALD_B                // LD.B rd, rj, si12
	ALD_BU               // LD.BU rd, rj, si12
	ALD_D                // LD.D rd, rj, si12
	ALD_H                // LD.H rd, rj, si12
	ALD_HU               // LD.HU rd, rj, si12
	ALD_W                // LD.W rd, rj, si12
	ALD_WU               // LD.WU rd, rj, si12
	ALLACQ_D             // LLACQ.D rd, rj
	ALLACQ_W             // LLACQ.W rd, rj
	ALL_D                // LL.D rd, rj, si14
	ALL_W                // LL.W rd, rj, si14
	ALU12I_W             // LU12I.W rd, si20
	ALU32I_D             // LU32I.D rd, si20
	ALU52I_D             // LU52I.D rd, rj, si12
	AMASKEQZ             // MASKEQZ rd, rj, rk
	AMASKNEZ             // MASKNEZ rd, rj, rk
	AMOD_D               // MOD.D rd, rj, rk
	AMOD_DU              // MOD.DU rd, rj, rk
	AMOD_W               // MOD.W rd, rj, rk
	AMOD_WU              // MOD.WU rd, rj, rk
	AMOVCF2FR            // MOVCF2FR fd, cj
	AMOVCF2GR            // MOVCF2GR rd, cj
	AMOVFCSR2GR          // MOVFCSR2GR rd, fcsr
	AMOVFR2CF            // MOVFR2CF cd, fj
	AMOVFR2GR_D          // MOVFR2GR.D rd, fj
	AMOVFR2GR_S          // MOVFR2GR.S rd, fj
	AMOVFRH2GR_S         // MOVFRH2GR.S rd, fj
	AMOVGR2CF            // MOVGR2CF cd, rj
	AMOVGR2FCSR          // MOVGR2FCSR fcsr, rj
	AMOVGR2FRH_W         // MOVGR2FRH.W fd, rj
	AMOVGR2FR_D          // MOVGR2FR.D fd, rj
	AMOVGR2FR_W          // MOVGR2FR.W fd, rj
	AMULH_D              // MULH.D rd, rj, rk
	AMULH_DU             // MULH.DU rd, rj, rk
	AMULH_W              // MULH.W rd, rj, rk
	AMULH_WU             // MULH.WU rd, rj, rk
	AMULW_D_W            // MULW.D.W rd, rj, rk
	AMULW_D_WU           // MULW.D.WU rd, rj, rk
	AMUL_D               // MUL.D rd, rj, rk
	AMUL_W               // MUL.W rd, rj, rk
	ANOR                 // NOR rd, rj, rk
	AOR                  // OR rd, rj, rk
	AORI                 // ORI rd, rj, ui12
	AORN                 // ORN rd, rj, rk
	APCADDI              // PCADDI rd, si20
	APCADDU12I           // PCADDU12I rd, si20
	APCADDU18I           // PCADDU18I rd, si20
	APCALAU12I           // PCALAU12I rd, si20
	APRELD               // PRELD hint, rj, si12
	APRELDX              // PRELDX hint, rj, rk
	ARDTIMEH_W           // RDTIMEH.W rd, rj
	ARDTIMEL_W           // RDTIMEL.W rd, rj
	ARDTIME_D            // RDTIME.D rd, rj
	AREVB_2H             // REVB.2H rd, rj
	AREVB_2W             // REVB.2W rd, rj
	AREVB_4H             // REVB.4H rd, rj
	AREVB_D              // REVB.D rd, rj
	AREVH_2W             // REVH.2W rd, rj
	AREVH_D              // REVH.D rd, rj
	AROTRI_D             // ROTRI.D rd, rj, ui6
	AROTRI_W             // ROTRI.W rd, rj, ui5
	AROTR_D              // ROTR.D rd, rj, rk
	AROTR_W              // ROTR.W rd, rj, rk
	ASCREL_D             // SCREL.D rd, rj
	ASCREL_W             // SCREL.W rd, rj
	ASC_D                // SC.D rd, rj, si14
	ASC_Q                // SC.Q rd, rk, rj
	ASC_W                // SC.W rd, rj, si14
	ASLLI_D              // SLLI.D rd, rj, ui6
	ASLLI_W              // SLLI.W rd, rj, ui5
	ASLL_D               // SLL.D rd, rj, rk
	ASLL_W               // SLL.W rd, rj, rk
	ASLT                 // SLT rd, rj, rk
	ASLTI                // SLTI rd, rj, si12
	ASLTU                // SLTU rd, rj, rk
	ASLTUI               // SLTUI rd, rj, si12
	ASRAI_D              // SRAI.D rd, rj, ui6
	ASRAI_W              // SRAI.W rd, rj, ui5
	ASRA_D               // SRA.D rd, rj, rk
	ASRA_W               // SRA.W rd, rj, rk
	ASRLI_D              // SRLI.D rd, rj, ui6
	ASRLI_W              // SRLI.W rd, rj, ui5
	ASRL_D               // SRL.D rd, rj, rk
	ASRL_W               // SRL.W rd, rj, rk
	ASTGT_B              // STGT.B rd, rj, rk
	ASTGT_D              // STGT.D rd, rj, rk
	ASTGT_H              // STGT.H rd, rj, rk
	ASTGT_W              // STGT.W rd, rj, rk
	ASTLE_B              // STLE.B rd, rj, rk
	ASTLE_D              // STLE.D rd, rj, rk
	ASTLE_H              // STLE.H rd, rj, rk
	ASTLE_W              // STLE.W rd, rj, rk
	ASTPTR_D             // STPTR.D rd, rj, si14
	ASTPTR_W             // STPTR.W rd, rj, si14
	ASTX_B               // STX.B rd, rj, rk
	ASTX_D               // STX.D rd, rj, rk
	ASTX_H               // STX.H rd, rj, rk
	ASTX_W               // STX.W rd, rj, rk
	AST_B                // ST.B rd, rj, si12
	AST_D                // ST.D rd, rj, si12
	AST_H                // ST.H rd, rj, si12
	AST_W                // ST.W rd, rj, si12
	ASUB_D               // SUB.D rd, rj, rk
	ASUB_W               // SUB.W rd, rj, rk
	ASYSCALL             // SYSCALL code
	ATLBCLR              // TLBCLR
	ATLBFILL             // TLBFILL
	ATLBFLUSH            // TLBFLUSH
	ATLBRD               // TLBRD
	ATLBSRCH             // TLBSRCH
	ATLBWR               // TLBWR
	AXOR                 // XOR rd, rj, rk
	AXORI                // XORI rd, rj, ui12

	ALAST
)

var _Anames = [...]string{
	AADDI_D:       "ADDI.D",
	AADDI_W:       "ADDI.W",
	AADDU16I_D:    "ADDU16I.D",
	AADD_D:        "ADD.D",
	AADD_W:        "ADD.W",
	AALSL_D:       "ALSL.D",
	AALSL_W:       "ALSL.W",
	AALSL_WU:      "ALSL.WU",
	AAMADD_B:      "AMADD.B",
	AAMADD_D:      "AMADD.D",
	AAMADD_DB_B:   "AMADD_DB.B",
	AAMADD_DB_D:   "AMADD_DB.D",
	AAMADD_DB_H:   "AMADD_DB.H",
	AAMADD_DB_W:   "AMADD_DB.W",
	AAMADD_H:      "AMADD.H",
	AAMADD_W:      "AMADD.W",
	AAMAND_D:      "AMAND.D",
	AAMAND_DB_D:   "AMAND_DB.D",
	AAMAND_DB_W:   "AMAND_DB.W",
	AAMAND_W:      "AMAND.W",
	AAMCAS_B:      "AMCAS.B",
	AAMCAS_D:      "AMCAS.D",
	AAMCAS_DB_B:   "AMCAS_DB.B",
	AAMCAS_DB_D:   "AMCAS_DB.D",
	AAMCAS_DB_H:   "AMCAS_DB.H",
	AAMCAS_DB_W:   "AMCAS_DB.W",
	AAMCAS_H:      "AMCAS.H",
	AAMCAS_W:      "AMCAS.W",
	AAMMAX_D:      "AMMAX.D",
	AAMMAX_DB_D:   "AMMAX_DB.D",
	AAMMAX_DB_DU:  "AMMAX_DB.DU",
	AAMMAX_DB_W:   "AMMAX_DB.W",
	AAMMAX_DB_WU:  "AMMAX_DB.WU",
	AAMMAX_DU:     "AMMAX.DU",
	AAMMAX_W:      "AMMAX.W",
	AAMMAX_WU:     "AMMAX.WU",
	AAMMIN_D:      "AMMIN.D",
	AAMMIN_DB_D:   "AMMIN_DB.D",
	AAMMIN_DB_DU:  "AMMIN_DB.DU",
	AAMMIN_DB_W:   "AMMIN_DB.W",
	AAMMIN_DB_WU:  "AMMIN_DB.WU",
	AAMMIN_DU:     "AMMIN.DU",
	AAMMIN_W:      "AMMIN.W",
	AAMMIN_WU:     "AMMIN.WU",
	AAMOR_D:       "AMOR.D",
	AAMOR_DB_D:    "AMOR_DB.D",
	AAMOR_DB_W:    "AMOR_DB.W",
	AAMOR_W:       "AMOR.W",
	AAMSWAP_B:     "AMSWAP.B",
	AAMSWAP_D:     "AMSWAP.D",
	AAMSWAP_DB_B:  "AMSWAP_DB.B",
	AAMSWAP_DB_D:  "AMSWAP_DB.D",
	AAMSWAP_DB_H:  "AMSWAP_DB.H",
	AAMSWAP_DB_W:  "AMSWAP_DB.W",
	AAMSWAP_H:     "AMSWAP.H",
	AAMSWAP_W:     "AMSWAP.W",
	AAMXOR_D:      "AMXOR.D",
	AAMXOR_DB_D:   "AMXOR_DB.D",
	AAMXOR_DB_W:   "AMXOR_DB.W",
	AAMXOR_W:      "AMXOR.W",
	AAND:          "AND",
	AANDI:         "ANDI",
	AANDN:         "ANDN",
	AASRTGT_D:     "ASRTGT.D",
	AASRTLE_D:     "ASRTLE.D",
	AB:            "B",
	ABCEQZ:        "BCEQZ",
	ABCNEZ:        "BCNEZ",
	ABEQ:          "BEQ",
	ABEQZ:         "BEQZ",
	ABGE:          "BGE",
	ABGEU:         "BGEU",
	ABITREV_4B:    "BITREV.4B",
	ABITREV_8B:    "BITREV.8B",
	ABITREV_D:     "BITREV.D",
	ABITREV_W:     "BITREV.W",
	ABL:           "BL",
	ABLT:          "BLT",
	ABLTU:         "BLTU",
	ABNE:          "BNE",
	ABNEZ:         "BNEZ",
	ABREAK:        "BREAK",
	ABSTRINS_D:    "BSTRINS.D",
	ABSTRINS_W:    "BSTRINS.W",
	ABSTRPICK_D:   "BSTRPICK.D",
	ABSTRPICK_W:   "BSTRPICK.W",
	ABYTEPICK_D:   "BYTEPICK.D",
	ABYTEPICK_W:   "BYTEPICK.W",
	ACACOP:        "CACOP",
	ACLO_D:        "CLO.D",
	ACLO_W:        "CLO.W",
	ACLZ_D:        "CLZ.D",
	ACLZ_W:        "CLZ.W",
	ACPUCFG:       "CPUCFG",
	ACRCC_W_B_W:   "CRCC.W.B.W",
	ACRCC_W_D_W:   "CRCC.W.D.W",
	ACRCC_W_H_W:   "CRCC.W.H.W",
	ACRCC_W_W_W:   "CRCC.W.W.W",
	ACRC_W_B_W:    "CRC.W.B.W",
	ACRC_W_D_W:    "CRC.W.D.W",
	ACRC_W_H_W:    "CRC.W.H.W",
	ACRC_W_W_W:    "CRC.W.W.W",
	ACSRRD:        "CSRRD",
	ACSRWR:        "CSRWR",
	ACSRXCHG:      "CSRXCHG",
	ACTO_D:        "CTO.D",
	ACTO_W:        "CTO.W",
	ACTZ_D:        "CTZ.D",
	ACTZ_W:        "CTZ.W",
	ADBAR:         "DBAR",
	ADBCL:         "DBCL",
	ADIV_D:        "DIV.D",
	ADIV_DU:       "DIV.DU",
	ADIV_W:        "DIV.W",
	ADIV_WU:       "DIV.WU",
	AERTN:         "ERTN",
	AEXT_W_B:      "EXT.W.B",
	AEXT_W_H:      "EXT.W.H",
	AFABS_D:       "FABS.D",
	AFABS_S:       "FABS.S",
	AFADD_D:       "FADD.D",
	AFADD_S:       "FADD.S",
	AFCLASS_D:     "FCLASS.D",
	AFCLASS_S:     "FCLASS.S",
	AFCMP_CAF_D:   "FCMP.CAF.D",
	AFCMP_CAF_S:   "FCMP.CAF.S",
	AFCMP_CEQ_D:   "FCMP.CEQ.D",
	AFCMP_CEQ_S:   "FCMP.CEQ.S",
	AFCMP_CLE_D:   "FCMP.CLE.D",
	AFCMP_CLE_S:   "FCMP.CLE.S",
	AFCMP_CLT_D:   "FCMP.CLT.D",
	AFCMP_CLT_S:   "FCMP.CLT.S",
	AFCMP_CNE_D:   "FCMP.CNE.D",
	AFCMP_CNE_S:   "FCMP.CNE.S",
	AFCMP_COR_D:   "FCMP.COR.D",
	AFCMP_COR_S:   "FCMP.COR.S",
	AFCMP_CUEQ_D:  "FCMP.CUEQ.D",
	AFCMP_CUEQ_S:  "FCMP.CUEQ.S",
	AFCMP_CULE_D:  "FCMP.CULE.D",
	AFCMP_CULE_S:  "FCMP.CULE.S",
	AFCMP_CULT_D:  "FCMP.CULT.D",
	AFCMP_CULT_S:  "FCMP.CULT.S",
	AFCMP_CUNE_D:  "FCMP.CUNE.D",
	AFCMP_CUNE_S:  "FCMP.CUNE.S",
	AFCMP_CUN_D:   "FCMP.CUN.D",
	AFCMP_CUN_S:   "FCMP.CUN.S",
	AFCMP_SAF_D:   "FCMP.SAF.D",
	AFCMP_SAF_S:   "FCMP.SAF.S",
	AFCMP_SEQ_D:   "FCMP.SEQ.D",
	AFCMP_SEQ_S:   "FCMP.SEQ.S",
	AFCMP_SLE_D:   "FCMP.SLE.D",
	AFCMP_SLE_S:   "FCMP.SLE.S",
	AFCMP_SLT_D:   "FCMP.SLT.D",
	AFCMP_SLT_S:   "FCMP.SLT.S",
	AFCMP_SNE_D:   "FCMP.SNE.D",
	AFCMP_SNE_S:   "FCMP.SNE.S",
	AFCMP_SOR_D:   "FCMP.SOR.D",
	AFCMP_SOR_S:   "FCMP.SOR.S",
	AFCMP_SUEQ_D:  "FCMP.SUEQ.D",
	AFCMP_SUEQ_S:  "FCMP.SUEQ.S",
	AFCMP_SULE_D:  "FCMP.SULE.D",
	AFCMP_SULE_S:  "FCMP.SULE.S",
	AFCMP_SULT_D:  "FCMP.SULT.D",
	AFCMP_SULT_S:  "FCMP.SULT.S",
	AFCMP_SUNE_D:  "FCMP.SUNE.D",
	AFCMP_SUNE_S:  "FCMP.SUNE.S",
	AFCMP_SUN_D:   "FCMP.SUN.D",
	AFCMP_SUN_S:   "FCMP.SUN.S",
	AFCOPYSIGN_D:  "FCOPYSIGN.D",
	AFCOPYSIGN_S:  "FCOPYSIGN.S",
	AFCVT_D_S:     "FCVT.D.S",
	AFCVT_S_D:     "FCVT.S.D",
	AFDIV_D:       "FDIV.D",
	AFDIV_S:       "FDIV.S",
	AFFINT_D_L:    "FFINT.D.L",
	AFFINT_D_W:    "FFINT.D.W",
	AFFINT_S_L:    "FFINT.S.L",
	AFFINT_S_W:    "FFINT.S.W",
	AFLDGT_D:      "FLDGT.D",
	AFLDGT_S:      "FLDGT.S",
	AFLDLE_D:      "FLDLE.D",
	AFLDLE_S:      "FLDLE.S",
	AFLDX_D:       "FLDX.D",
	AFLDX_S:       "FLDX.S",
	AFLD_D:        "FLD.D",
	AFLD_S:        "FLD.S",
	AFLOGB_D:      "FLOGB.D",
	AFLOGB_S:      "FLOGB.S",
	AFMADD_D:      "FMADD.D",
	AFMADD_S:      "FMADD.S",
	AFMAXA_D:      "FMAXA.D",
	AFMAXA_S:      "FMAXA.S",
	AFMAX_D:       "FMAX.D",
	AFMAX_S:       "FMAX.S",
	AFMINA_D:      "FMINA.D",
	AFMINA_S:      "FMINA.S",
	AFMIN_D:       "FMIN.D",
	AFMIN_S:       "FMIN.S",
	AFMOV_D:       "FMOV.D",
	AFMOV_S:       "FMOV.S",
	AFMSUB_D:      "FMSUB.D",
	AFMSUB_S:      "FMSUB.S",
	AFMUL_D:       "FMUL.D",
	AFMUL_S:       "FMUL.S",
	AFNEG_D:       "FNEG.D",
	AFNEG_S:       "FNEG.S",
	AFNMADD_D:     "FNMADD.D",
	AFNMADD_S:     "FNMADD.S",
	AFNMSUB_D:     "FNMSUB.D",
	AFNMSUB_S:     "FNMSUB.S",
	AFRECIPE_D:    "FRECIPE.D",
	AFRECIPE_S:    "FRECIPE.S",
	AFRECIP_D:     "FRECIP.D",
	AFRECIP_S:     "FRECIP.S",
	AFRINT_D:      "FRINT.D",
	AFRINT_S:      "FRINT.S",
	AFRSQRTE_D:    "FRSQRTE.D",
	AFRSQRTE_S:    "FRSQRTE.S",
	AFRSQRT_D:     "FRSQRT.D",
	AFRSQRT_S:     "FRSQRT.S",
	AFSCALEB_D:    "FSCALEB.D",
	AFSCALEB_S:    "FSCALEB.S",
	AFSEL:         "FSEL",
	AFSQRT_D:      "FSQRT.D",
	AFSQRT_S:      "FSQRT.S",
	AFSTGT_D:      "FSTGT.D",
	AFSTGT_S:      "FSTGT.S",
	AFSTLE_D:      "FSTLE.D",
	AFSTLE_S:      "FSTLE.S",
	AFSTX_D:       "FSTX.D",
	AFSTX_S:       "FSTX.S",
	AFST_D:        "FST.D",
	AFST_S:        "FST.S",
	AFSUB_D:       "FSUB.D",
	AFSUB_S:       "FSUB.S",
	AFTINTRM_L_D:  "FTINTRM.L.D",
	AFTINTRM_L_S:  "FTINTRM.L.S",
	AFTINTRM_W_D:  "FTINTRM.W.D",
	AFTINTRM_W_S:  "FTINTRM.W.S",
	AFTINTRNE_L_D: "FTINTRNE.L.D",
	AFTINTRNE_L_S: "FTINTRNE.L.S",
	AFTINTRNE_W_D: "FTINTRNE.W.D",
	AFTINTRNE_W_S: "FTINTRNE.W.S",
	AFTINTRP_L_D:  "FTINTRP.L.D",
	AFTINTRP_L_S:  "FTINTRP.L.S",
	AFTINTRP_W_D:  "FTINTRP.W.D",
	AFTINTRP_W_S:  "FTINTRP.W.S",
	AFTINTRZ_L_D:  "FTINTRZ.L.D",
	AFTINTRZ_L_S:  "FTINTRZ.L.S",
	AFTINTRZ_W_D:  "FTINTRZ.W.D",
	AFTINTRZ_W_S:  "FTINTRZ.W.S",
	AFTINT_L_D:    "FTINT.L.D",
	AFTINT_L_S:    "FTINT.L.S",
	AFTINT_W_D:    "FTINT.W.D",
	AFTINT_W_S:    "FTINT.W.S",
	AIBAR:         "IBAR",
	AIDLE:         "IDLE",
	AINVTLB:       "INVTLB",
	AIOCSRRD_B:    "IOCSRRD.B",
	AIOCSRRD_D:    "IOCSRRD.D",
	AIOCSRRD_H:    "IOCSRRD.H",
	AIOCSRRD_W:    "IOCSRRD.W",
	AIOCSRWR_B:    "IOCSRWR.B",
	AIOCSRWR_D:    "IOCSRWR.D",
	AIOCSRWR_H:    "IOCSRWR.H",
	AIOCSRWR_W:    "IOCSRWR.W",
	AJIRL:         "JIRL",
	ALDDIR:        "LDDIR",
	ALDGT_B:       "LDGT.B",
	ALDGT_D:       "LDGT.D",
	ALDGT_H:       "LDGT.H",
	ALDGT_W:       "LDGT.W",
	ALDLE_B:       "LDLE.B",
	ALDLE_D:       "LDLE.D",
	ALDLE_H:       "LDLE.H",
	ALDLE_W:       "LDLE.W",
	ALDPTE:        "LDPTE",
	ALDPTR_D:      "LDPTR.D",
	ALDPTR_W:      "LDPTR.W",
	ALDX_B:        "LDX.B",
	ALDX_BU:       "LDX.BU",
	ALDX_D:        "LDX.D",
	ALDX_H:        "LDX.H",
	ALDX_HU:       "LDX.HU",
	ALDX_W:        "LDX.W",
	ALDX_WU:       "LDX.WU",
	ALD_B:         "LD.B",
	ALD_BU:        "LD.BU",
	ALD_D:         "LD.D",
	ALD_H:         "LD.H",
	ALD_HU:        "LD.HU",
	ALD_W:         "LD.W",
	ALD_WU:        "LD.WU",
	ALLACQ_D:      "LLACQ.D",
	ALLACQ_W:      "LLACQ.W",
	ALL_D:         "LL.D",
	ALL_W:         "LL.W",
	ALU12I_W:      "LU12I.W",
	ALU32I_D:      "LU32I.D",
	ALU52I_D:      "LU52I.D",
	AMASKEQZ:      "MASKEQZ",
	AMASKNEZ:      "MASKNEZ",
	AMOD_D:        "MOD.D",
	AMOD_DU:       "MOD.DU",
	AMOD_W:        "MOD.W",
	AMOD_WU:       "MOD.WU",
	AMOVCF2FR:     "MOVCF2FR",
	AMOVCF2GR:     "MOVCF2GR",
	AMOVFCSR2GR:   "MOVFCSR2GR",
	AMOVFR2CF:     "MOVFR2CF",
	AMOVFR2GR_D:   "MOVFR2GR.D",
	AMOVFR2GR_S:   "MOVFR2GR.S",
	AMOVFRH2GR_S:  "MOVFRH2GR.S",
	AMOVGR2CF:     "MOVGR2CF",
	AMOVGR2FCSR:   "MOVGR2FCSR",
	AMOVGR2FRH_W:  "MOVGR2FRH.W",
	AMOVGR2FR_D:   "MOVGR2FR.D",
	AMOVGR2FR_W:   "MOVGR2FR.W",
	AMULH_D:       "MULH.D",
	AMULH_DU:      "MULH.DU",
	AMULH_W:       "MULH.W",
	AMULH_WU:      "MULH.WU",
	AMULW_D_W:     "MULW.D.W",
	AMULW_D_WU:    "MULW.D.WU",
	AMUL_D:        "MUL.D",
	AMUL_W:        "MUL.W",
	ANOR:          "NOR",
	AOR:           "OR",
	AORI:          "ORI",
	AORN:          "ORN",
	APCADDI:       "PCADDI",
	APCADDU12I:    "PCADDU12I",
	APCADDU18I:    "PCADDU18I",
	APCALAU12I:    "PCALAU12I",
	APRELD:        "PRELD",
	APRELDX:       "PRELDX",
	ARDTIMEH_W:    "RDTIMEH.W",
	ARDTIMEL_W:    "RDTIMEL.W",
	ARDTIME_D:     "RDTIME.D",
	AREVB_2H:      "REVB.2H",
	AREVB_2W:      "REVB.2W",
	AREVB_4H:      "REVB.4H",
	AREVB_D:       "REVB.D",
	AREVH_2W:      "REVH.2W",
	AREVH_D:       "REVH.D",
	AROTRI_D:      "ROTRI.D",
	AROTRI_W:      "ROTRI.W",
	AROTR_D:       "ROTR.D",
	AROTR_W:       "ROTR.W",
	ASCREL_D:      "SCREL.D",
	ASCREL_W:      "SCREL.W",
	ASC_D:         "SC.D",
	ASC_Q:         "SC.Q",
	ASC_W:         "SC.W",
	ASLLI_D:       "SLLI.D",
	ASLLI_W:       "SLLI.W",
	ASLL_D:        "SLL.D",
	ASLL_W:        "SLL.W",
	ASLT:          "SLT",
	ASLTI:         "SLTI",
	ASLTU:         "SLTU",
	ASLTUI:        "SLTUI",
	ASRAI_D:       "SRAI.D",
	ASRAI_W:       "SRAI.W",
	ASRA_D:        "SRA.D",
	ASRA_W:        "SRA.W",
	ASRLI_D:       "SRLI.D",
	ASRLI_W:       "SRLI.W",
	ASRL_D:        "SRL.D",
	ASRL_W:        "SRL.W",
	ASTGT_B:       "STGT.B",
	ASTGT_D:       "STGT.D",
	ASTGT_H:       "STGT.H",
	ASTGT_W:       "STGT.W",
	ASTLE_B:       "STLE.B",
	ASTLE_D:       "STLE.D",
	ASTLE_H:       "STLE.H",
	ASTLE_W:       "STLE.W",
	ASTPTR_D:      "STPTR.D",
	ASTPTR_W:      "STPTR.W",
	ASTX_B:        "STX.B",
	ASTX_D:        "STX.D",
	ASTX_H:        "STX.H",
	ASTX_W:        "STX.W",
	AST_B:         "ST.B",
	AST_D:         "ST.D",
	AST_H:         "ST.H",
	AST_W:         "ST.W",
	ASUB_D:        "SUB.D",
	ASUB_W:        "SUB.W",
	ASYSCALL:      "SYSCALL",
	ATLBCLR:       "TLBCLR",
	ATLBFILL:      "TLBFILL",
	ATLBFLUSH:     "TLBFLUSH",
	ATLBRD:        "TLBRD",
	ATLBSRCH:      "TLBSRCH",
	ATLBWR:        "TLBWR",
	AXOR:          "XOR",
	AXORI:         "XORI",
}

// 指令编码信息表
var _AOpContextTable = [...]_OpContextType{
	AADDI_D:       {mask: 0xffc00000, value: 0x02c00000, op: AADDI_D, fmt: OpFormatType_2R_si12},
	AADDI_W:       {mask: 0xffc00000, value: 0x02800000, op: AADDI_W, fmt: OpFormatType_2R_si12},
	AADDU16I_D:    {mask: 0xfc000000, value: 0x10000000, op: AADDU16I_D, fmt: OpFormatType_2R_si14},
	AADD_D:        {mask: 0xffff8000, value: 0x00108000, op: AADD_D, fmt: OpFormatType_3R},
	AADD_W:        {mask: 0xffff8000, value: 0x00100000, op: AADD_W, fmt: OpFormatType_3R},
	AALSL_D:       {mask: 0xfffe0000, value: 0x002c0000, op: AALSL_D, fmt: OpFormatType_3R_s2},
	AALSL_W:       {mask: 0xfffe0000, value: 0x00040000, op: AALSL_W, fmt: OpFormatType_3R_s2},
	AALSL_WU:      {mask: 0xfffe0000, value: 0x00060000, op: AALSL_WU, fmt: OpFormatType_3R_s2},
	AAMADD_B:      {mask: 0xffff8000, value: 0x385d0000, op: AAMADD_B, fmt: OpFormatType_3R},
	AAMADD_D:      {mask: 0xffff8000, value: 0x38618000, op: AAMADD_D, fmt: OpFormatType_3R},
	AAMADD_DB_B:   {mask: 0xffff8000, value: 0x385f0000, op: AAMADD_DB_B, fmt: OpFormatType_3R},
	AAMADD_DB_D:   {mask: 0xffff8000, value: 0x386a8000, op: AAMADD_DB_D, fmt: OpFormatType_3R},
	AAMADD_DB_H:   {mask: 0xffff8000, value: 0x385f8000, op: AAMADD_DB_H, fmt: OpFormatType_3R},
	AAMADD_DB_W:   {mask: 0xffff8000, value: 0x386a0000, op: AAMADD_DB_W, fmt: OpFormatType_3R},
	AAMADD_H:      {mask: 0xffff8000, value: 0x385d8000, op: AAMADD_H, fmt: OpFormatType_3R},
	AAMADD_W:      {mask: 0xffff8000, value: 0x38610000, op: AAMADD_W, fmt: OpFormatType_3R},
	AAMAND_D:      {mask: 0xffff8000, value: 0x38628000, op: AAMAND_D, fmt: OpFormatType_3R},
	AAMAND_DB_D:   {mask: 0xffff8000, value: 0x386b8000, op: AAMAND_DB_D, fmt: OpFormatType_3R},
	AAMAND_DB_W:   {mask: 0xffff8000, value: 0x386b0000, op: AAMAND_DB_W, fmt: OpFormatType_3R},
	AAMAND_W:      {mask: 0xffff8000, value: 0x38620000, op: AAMAND_W, fmt: OpFormatType_3R},
	AAMCAS_B:      {mask: 0xffff8000, value: 0x38580000, op: AAMCAS_B, fmt: OpFormatType_3R},
	AAMCAS_D:      {mask: 0xffff8000, value: 0x38598000, op: AAMCAS_D, fmt: OpFormatType_3R},
	AAMCAS_DB_B:   {mask: 0xffff8000, value: 0x385a0000, op: AAMCAS_DB_B, fmt: OpFormatType_3R},
	AAMCAS_DB_D:   {mask: 0xffff8000, value: 0x385b8000, op: AAMCAS_DB_D, fmt: OpFormatType_3R},
	AAMCAS_DB_H:   {mask: 0xffff8000, value: 0x385a8000, op: AAMCAS_DB_H, fmt: OpFormatType_3R},
	AAMCAS_DB_W:   {mask: 0xffff8000, value: 0x385b0000, op: AAMCAS_DB_W, fmt: OpFormatType_3R},
	AAMCAS_H:      {mask: 0xffff8000, value: 0x38588000, op: AAMCAS_H, fmt: OpFormatType_3R},
	AAMCAS_W:      {mask: 0xffff8000, value: 0x38590000, op: AAMCAS_W, fmt: OpFormatType_3R},
	AAMMAX_D:      {mask: 0xffff8000, value: 0x38658000, op: AAMMAX_D, fmt: OpFormatType_3R},
	AAMMAX_DB_D:   {mask: 0xffff8000, value: 0x386e8000, op: AAMMAX_DB_D, fmt: OpFormatType_3R},
	AAMMAX_DB_DU:  {mask: 0xffff8000, value: 0x38708000, op: AAMMAX_DB_DU, fmt: OpFormatType_3R},
	AAMMAX_DB_W:   {mask: 0xffff8000, value: 0x386e0000, op: AAMMAX_DB_W, fmt: OpFormatType_3R},
	AAMMAX_DB_WU:  {mask: 0xffff8000, value: 0x38700000, op: AAMMAX_DB_WU, fmt: OpFormatType_3R},
	AAMMAX_DU:     {mask: 0xffff8000, value: 0x38678000, op: AAMMAX_DU, fmt: OpFormatType_3R},
	AAMMAX_W:      {mask: 0xffff8000, value: 0x38650000, op: AAMMAX_W, fmt: OpFormatType_3R},
	AAMMAX_WU:     {mask: 0xffff8000, value: 0x38670000, op: AAMMAX_WU, fmt: OpFormatType_3R},
	AAMMIN_D:      {mask: 0xffff8000, value: 0x38668000, op: AAMMIN_D, fmt: OpFormatType_3R},
	AAMMIN_DB_D:   {mask: 0xffff8000, value: 0x386f8000, op: AAMMIN_DB_D, fmt: OpFormatType_3R},
	AAMMIN_DB_DU:  {mask: 0xffff8000, value: 0x38718000, op: AAMMIN_DB_DU, fmt: OpFormatType_3R},
	AAMMIN_DB_W:   {mask: 0xffff8000, value: 0x386f0000, op: AAMMIN_DB_W, fmt: OpFormatType_3R},
	AAMMIN_DB_WU:  {mask: 0xffff8000, value: 0x38710000, op: AAMMIN_DB_WU, fmt: OpFormatType_3R},
	AAMMIN_DU:     {mask: 0xffff8000, value: 0x38688000, op: AAMMIN_DU, fmt: OpFormatType_3R},
	AAMMIN_W:      {mask: 0xffff8000, value: 0x38660000, op: AAMMIN_W, fmt: OpFormatType_3R},
	AAMMIN_WU:     {mask: 0xffff8000, value: 0x38680000, op: AAMMIN_WU, fmt: OpFormatType_3R},
	AAMOR_D:       {mask: 0xffff8000, value: 0x38638000, op: AAMOR_D, fmt: OpFormatType_3R},
	AAMOR_DB_D:    {mask: 0xffff8000, value: 0x386c8000, op: AAMOR_DB_D, fmt: OpFormatType_3R},
	AAMOR_DB_W:    {mask: 0xffff8000, value: 0x386c0000, op: AAMOR_DB_W, fmt: OpFormatType_3R},
	AAMOR_W:       {mask: 0xffff8000, value: 0x38630000, op: AAMOR_W, fmt: OpFormatType_3R},
	AAMSWAP_B:     {mask: 0xffff8000, value: 0x385c0000, op: AAMSWAP_B, fmt: OpFormatType_3R},
	AAMSWAP_D:     {mask: 0xffff8000, value: 0x38608000, op: AAMSWAP_D, fmt: OpFormatType_3R},
	AAMSWAP_DB_B:  {mask: 0xffff8000, value: 0x385e0000, op: AAMSWAP_DB_B, fmt: OpFormatType_3R},
	AAMSWAP_DB_D:  {mask: 0xffff8000, value: 0x38698000, op: AAMSWAP_DB_D, fmt: OpFormatType_3R},
	AAMSWAP_DB_H:  {mask: 0xffff8000, value: 0x385e8000, op: AAMSWAP_DB_H, fmt: OpFormatType_3R},
	AAMSWAP_DB_W:  {mask: 0xffff8000, value: 0x38690000, op: AAMSWAP_DB_W, fmt: OpFormatType_3R},
	AAMSWAP_H:     {mask: 0xffff8000, value: 0x385c8000, op: AAMSWAP_H, fmt: OpFormatType_3R},
	AAMSWAP_W:     {mask: 0xffff8000, value: 0x38600000, op: AAMSWAP_W, fmt: OpFormatType_3R},
	AAMXOR_D:      {mask: 0xffff8000, value: 0x38648000, op: AAMXOR_D, fmt: OpFormatType_3R},
	AAMXOR_DB_D:   {mask: 0xffff8000, value: 0x386d8000, op: AAMXOR_DB_D, fmt: OpFormatType_3R},
	AAMXOR_DB_W:   {mask: 0xffff8000, value: 0x386d0000, op: AAMXOR_DB_W, fmt: OpFormatType_3R},
	AAMXOR_W:      {mask: 0xffff8000, value: 0x38640000, op: AAMXOR_W, fmt: OpFormatType_3R},
	AAND:          {mask: 0xffff8000, value: 0x00148000, op: AAND, fmt: OpFormatType_3R},
	AANDI:         {mask: 0xffc00000, value: 0x03400000, op: AANDI, fmt: OpFormatType_2R_ui12},
	AANDN:         {mask: 0xffff8000, value: 0x00168000, op: AANDN, fmt: OpFormatType_3R},
	AASRTGT_D:     {mask: 0xffff801f, value: 0x00018000, op: AASRTGT_D, fmt: OpFormatType_0_2R},
	AASRTLE_D:     {mask: 0xffff801f, value: 0x00010000, op: AASRTLE_D, fmt: OpFormatType_0_2R},
	AB:            {mask: 0xfc000000, value: 0x50000000, op: AB, fmt: OpFormatType_offset},
	ABCEQZ:        {mask: 0xfc000300, value: 0x48000000, op: ABCEQZ, fmt: OpFormatType_cj_offset},
	ABCNEZ:        {mask: 0xfc000300, value: 0x48000100, op: ABCNEZ, fmt: OpFormatType_cj_offset},
	ABEQ:          {mask: 0xfc000000, value: 0x58000000, op: ABEQ, fmt: OpFormatType_rj_rd_offset},
	ABEQZ:         {mask: 0xfc000000, value: 0x40000000, op: ABEQZ, fmt: OpFormatType_rj_offset},
	ABGE:          {mask: 0xfc000000, value: 0x64000000, op: ABGE, fmt: OpFormatType_rj_rd_offset},
	ABGEU:         {mask: 0xfc000000, value: 0x6c000000, op: ABGEU, fmt: OpFormatType_rj_rd_offset},
	ABITREV_4B:    {mask: 0xfffffc00, value: 0x00004800, op: ABITREV_4B, fmt: OpFormatType_2R},
	ABITREV_8B:    {mask: 0xfffffc00, value: 0x00004c00, op: ABITREV_8B, fmt: OpFormatType_2R},
	ABITREV_D:     {mask: 0xfffffc00, value: 0x00005400, op: ABITREV_D, fmt: OpFormatType_2R},
	ABITREV_W:     {mask: 0xfffffc00, value: 0x00005000, op: ABITREV_W, fmt: OpFormatType_2R},
	ABL:           {mask: 0xfc000000, value: 0x54000000, op: ABL, fmt: OpFormatType_offset},
	ABLT:          {mask: 0xfc000000, value: 0x60000000, op: ABLT, fmt: OpFormatType_rj_rd_offset},
	ABLTU:         {mask: 0xfc000000, value: 0x68000000, op: ABLTU, fmt: OpFormatType_rj_rd_offset},
	ABNE:          {mask: 0xfc000000, value: 0x5c000000, op: ABNE, fmt: OpFormatType_rj_rd_offset},
	ABNEZ:         {mask: 0xfc000000, value: 0x44000000, op: ABNEZ, fmt: OpFormatType_rj_offset},
	ABREAK:        {mask: 0xffff8000, value: 0x002a0000, op: ABREAK, fmt: OpFormatType_code},
	ABSTRINS_D:    {mask: 0xffc00000, value: 0x00800000, op: ABSTRINS_D, fmt: OpFormatType_msbd_lsbd},
	ABSTRINS_W:    {mask: 0xffe08000, value: 0x00600000, op: ABSTRINS_W, fmt: OpFormatType_msbw_lsbw},
	ABSTRPICK_D:   {mask: 0xffc00000, value: 0x00c00000, op: ABSTRPICK_D, fmt: OpFormatType_msbd_lsbd},
	ABSTRPICK_W:   {mask: 0xffe08000, value: 0x00608000, op: ABSTRPICK_W, fmt: OpFormatType_msbw_lsbw},
	ABYTEPICK_D:   {mask: 0xfffc0000, value: 0x000c0000, op: ABYTEPICK_D, fmt: OpFormatType_3R_s3},
	ABYTEPICK_W:   {mask: 0xfffe0000, value: 0x00080000, op: ABYTEPICK_W, fmt: OpFormatType_3R_s2},
	ACACOP:        {mask: 0xffc00000, value: 0x06000000, op: ACACOP, fmt: OpFormatType_code_1R_si12},
	ACLO_D:        {mask: 0xfffffc00, value: 0x00002000, op: ACLO_D, fmt: OpFormatType_2R},
	ACLO_W:        {mask: 0xfffffc00, value: 0x00001000, op: ACLO_W, fmt: OpFormatType_2R},
	ACLZ_D:        {mask: 0xfffffc00, value: 0x00002400, op: ACLZ_D, fmt: OpFormatType_2R},
	ACLZ_W:        {mask: 0xfffffc00, value: 0x00001400, op: ACLZ_W, fmt: OpFormatType_2R},
	ACPUCFG:       {mask: 0xfffffc00, value: 0x00006c00, op: ACPUCFG, fmt: OpFormatType_2R},
	ACRCC_W_B_W:   {mask: 0xffff8000, value: 0x00260000, op: ACRCC_W_B_W, fmt: OpFormatType_3R},
	ACRCC_W_D_W:   {mask: 0xffff8000, value: 0x00278000, op: ACRCC_W_D_W, fmt: OpFormatType_3R},
	ACRCC_W_H_W:   {mask: 0xffff8000, value: 0x00268000, op: ACRCC_W_H_W, fmt: OpFormatType_3R},
	ACRCC_W_W_W:   {mask: 0xffff8000, value: 0x00270000, op: ACRCC_W_W_W, fmt: OpFormatType_3R},
	ACRC_W_B_W:    {mask: 0xffff8000, value: 0x00240000, op: ACRC_W_B_W, fmt: OpFormatType_3R},
	ACRC_W_D_W:    {mask: 0xffff8000, value: 0x00258000, op: ACRC_W_D_W, fmt: OpFormatType_3R},
	ACRC_W_H_W:    {mask: 0xffff8000, value: 0x00248000, op: ACRC_W_H_W, fmt: OpFormatType_3R},
	ACRC_W_W_W:    {mask: 0xffff8000, value: 0x00250000, op: ACRC_W_W_W, fmt: OpFormatType_3R},
	ACSRRD:        {mask: 0xff0003e0, value: 0x04000000, op: ACSRRD, fmt: OpFormatType_1R_csr},
	ACSRWR:        {mask: 0xff0003e0, value: 0x04000020, op: ACSRWR, fmt: OpFormatType_1R_csr},
	ACSRXCHG:      {mask: 0xff000000, value: 0x04000000, op: ACSRXCHG, fmt: OpFormatType_2R_csr},
	ACTO_D:        {mask: 0xfffffc00, value: 0x00002800, op: ACTO_D, fmt: OpFormatType_2R},
	ACTO_W:        {mask: 0xfffffc00, value: 0x00001800, op: ACTO_W, fmt: OpFormatType_2R},
	ACTZ_D:        {mask: 0xfffffc00, value: 0x00002c00, op: ACTZ_D, fmt: OpFormatType_2R},
	ACTZ_W:        {mask: 0xfffffc00, value: 0x00001c00, op: ACTZ_W, fmt: OpFormatType_2R},
	ADBAR:         {mask: 0xffff8000, value: 0x38720000, op: ADBAR, fmt: OpFormatType_hint},
	ADBCL:         {mask: 0xffff8000, value: 0x002a8000, op: ADBCL, fmt: OpFormatType_code},
	ADIV_D:        {mask: 0xffff8000, value: 0x00220000, op: ADIV_D, fmt: OpFormatType_3R},
	ADIV_DU:       {mask: 0xffff8000, value: 0x00230000, op: ADIV_DU, fmt: OpFormatType_3R},
	ADIV_W:        {mask: 0xffff8000, value: 0x00200000, op: ADIV_W, fmt: OpFormatType_3R},
	ADIV_WU:       {mask: 0xffff8000, value: 0x00210000, op: ADIV_WU, fmt: OpFormatType_3R},
	AERTN:         {mask: 0xffffffff, value: 0x06483800, op: AERTN, fmt: OpFormatType_NULL},
	AEXT_W_B:      {mask: 0xfffffc00, value: 0x00005c00, op: AEXT_W_B, fmt: OpFormatType_2R},
	AEXT_W_H:      {mask: 0xfffffc00, value: 0x00005800, op: AEXT_W_H, fmt: OpFormatType_2R},
	AFABS_D:       {mask: 0xfffffc00, value: 0x01140800, op: AFABS_D, fmt: OpFormatType_2F},
	AFABS_S:       {mask: 0xfffffc00, value: 0x01140400, op: AFABS_S, fmt: OpFormatType_2F},
	AFADD_D:       {mask: 0xffff8000, value: 0x01010000, op: AFADD_D, fmt: OpFormatType_3F},
	AFADD_S:       {mask: 0xffff8000, value: 0x01008000, op: AFADD_S, fmt: OpFormatType_3F},
	AFCLASS_D:     {mask: 0xfffffc00, value: 0x01143800, op: AFCLASS_D, fmt: OpFormatType_2F},
	AFCLASS_S:     {mask: 0xfffffc00, value: 0x01143400, op: AFCLASS_S, fmt: OpFormatType_2F},
	AFCMP_CAF_D:   {mask: 0xffff8018, value: 0x0c200000, op: AFCMP_CAF_D, fmt: OpFormatType_cd_2R},
	AFCMP_CAF_S:   {mask: 0xffff8018, value: 0x0c100000, op: AFCMP_CAF_S, fmt: OpFormatType_cd_2R},
	AFCMP_CEQ_D:   {mask: 0xffff8018, value: 0x0c220000, op: AFCMP_CEQ_D, fmt: OpFormatType_cd_2R},
	AFCMP_CEQ_S:   {mask: 0xffff8018, value: 0x0c120000, op: AFCMP_CEQ_S, fmt: OpFormatType_cd_2R},
	AFCMP_CLE_D:   {mask: 0xffff8018, value: 0x0c230000, op: AFCMP_CLE_D, fmt: OpFormatType_cd_2R},
	AFCMP_CLE_S:   {mask: 0xffff8018, value: 0x0c130000, op: AFCMP_CLE_S, fmt: OpFormatType_cd_2R},
	AFCMP_CLT_D:   {mask: 0xffff8018, value: 0x0c210000, op: AFCMP_CLT_D, fmt: OpFormatType_cd_2R},
	AFCMP_CLT_S:   {mask: 0xffff8018, value: 0x0c110000, op: AFCMP_CLT_S, fmt: OpFormatType_cd_2R},
	AFCMP_CNE_D:   {mask: 0xffff8018, value: 0x0c280000, op: AFCMP_CNE_D, fmt: OpFormatType_cd_2R},
	AFCMP_CNE_S:   {mask: 0xffff8018, value: 0x0c180000, op: AFCMP_CNE_S, fmt: OpFormatType_cd_2R},
	AFCMP_COR_D:   {mask: 0xffff8018, value: 0x0c2a0000, op: AFCMP_COR_D, fmt: OpFormatType_cd_2R},
	AFCMP_COR_S:   {mask: 0xffff8018, value: 0x0c1a0000, op: AFCMP_COR_S, fmt: OpFormatType_cd_2R},
	AFCMP_CUEQ_D:  {mask: 0xffff8018, value: 0x0c260000, op: AFCMP_CUEQ_D, fmt: OpFormatType_cd_2R},
	AFCMP_CUEQ_S:  {mask: 0xffff8018, value: 0x0c160000, op: AFCMP_CUEQ_S, fmt: OpFormatType_cd_2R},
	AFCMP_CULE_D:  {mask: 0xffff8018, value: 0x0c270000, op: AFCMP_CULE_D, fmt: OpFormatType_cd_2R},
	AFCMP_CULE_S:  {mask: 0xffff8018, value: 0x0c170000, op: AFCMP_CULE_S, fmt: OpFormatType_cd_2R},
	AFCMP_CULT_D:  {mask: 0xffff8018, value: 0x0c250000, op: AFCMP_CULT_D, fmt: OpFormatType_cd_2R},
	AFCMP_CULT_S:  {mask: 0xffff8018, value: 0x0c150000, op: AFCMP_CULT_S, fmt: OpFormatType_cd_2R},
	AFCMP_CUNE_D:  {mask: 0xffff8018, value: 0x0c2c0000, op: AFCMP_CUNE_D, fmt: OpFormatType_cd_2R},
	AFCMP_CUNE_S:  {mask: 0xffff8018, value: 0x0c1c0000, op: AFCMP_CUNE_S, fmt: OpFormatType_cd_2R},
	AFCMP_CUN_D:   {mask: 0xffff8018, value: 0x0c240000, op: AFCMP_CUN_D, fmt: OpFormatType_cd_2R},
	AFCMP_CUN_S:   {mask: 0xffff8018, value: 0x0c140000, op: AFCMP_CUN_S, fmt: OpFormatType_cd_2R},
	AFCMP_SAF_D:   {mask: 0xffff8018, value: 0x0c208000, op: AFCMP_SAF_D, fmt: OpFormatType_cd_2R},
	AFCMP_SAF_S:   {mask: 0xffff8018, value: 0x0c108000, op: AFCMP_SAF_S, fmt: OpFormatType_cd_2R},
	AFCMP_SEQ_D:   {mask: 0xffff8018, value: 0x0c228000, op: AFCMP_SEQ_D, fmt: OpFormatType_cd_2R},
	AFCMP_SEQ_S:   {mask: 0xffff8018, value: 0x0c128000, op: AFCMP_SEQ_S, fmt: OpFormatType_cd_2R},
	AFCMP_SLE_D:   {mask: 0xffff8018, value: 0x0c238000, op: AFCMP_SLE_D, fmt: OpFormatType_cd_2R},
	AFCMP_SLE_S:   {mask: 0xffff8018, value: 0x0c138000, op: AFCMP_SLE_S, fmt: OpFormatType_cd_2R},
	AFCMP_SLT_D:   {mask: 0xffff8018, value: 0x0c218000, op: AFCMP_SLT_D, fmt: OpFormatType_cd_2R},
	AFCMP_SLT_S:   {mask: 0xffff8018, value: 0x0c118000, op: AFCMP_SLT_S, fmt: OpFormatType_cd_2R},
	AFCMP_SNE_D:   {mask: 0xffff8018, value: 0x0c288000, op: AFCMP_SNE_D, fmt: OpFormatType_cd_2R},
	AFCMP_SNE_S:   {mask: 0xffff8018, value: 0x0c188000, op: AFCMP_SNE_S, fmt: OpFormatType_cd_2R},
	AFCMP_SOR_D:   {mask: 0xffff8018, value: 0x0c2a8000, op: AFCMP_SOR_D, fmt: OpFormatType_cd_2R},
	AFCMP_SOR_S:   {mask: 0xffff8018, value: 0x0c1a8000, op: AFCMP_SOR_S, fmt: OpFormatType_cd_2R},
	AFCMP_SUEQ_D:  {mask: 0xffff8018, value: 0x0c268000, op: AFCMP_SUEQ_D, fmt: OpFormatType_cd_2R},
	AFCMP_SUEQ_S:  {mask: 0xffff8018, value: 0x0c168000, op: AFCMP_SUEQ_S, fmt: OpFormatType_cd_2R},
	AFCMP_SULE_D:  {mask: 0xffff8018, value: 0x0c278000, op: AFCMP_SULE_D, fmt: OpFormatType_cd_2R},
	AFCMP_SULE_S:  {mask: 0xffff8018, value: 0x0c178000, op: AFCMP_SULE_S, fmt: OpFormatType_cd_2R},
	AFCMP_SULT_D:  {mask: 0xffff8018, value: 0x0c258000, op: AFCMP_SULT_D, fmt: OpFormatType_cd_2R},
	AFCMP_SULT_S:  {mask: 0xffff8018, value: 0x0c158000, op: AFCMP_SULT_S, fmt: OpFormatType_cd_2R},
	AFCMP_SUNE_D:  {mask: 0xffff8018, value: 0x0c2c8000, op: AFCMP_SUNE_D, fmt: OpFormatType_cd_2R},
	AFCMP_SUNE_S:  {mask: 0xffff8018, value: 0x0c1c8000, op: AFCMP_SUNE_S, fmt: OpFormatType_cd_2R},
	AFCMP_SUN_D:   {mask: 0xffff8018, value: 0x0c248000, op: AFCMP_SUN_D, fmt: OpFormatType_cd_2R},
	AFCMP_SUN_S:   {mask: 0xffff8018, value: 0x0c148000, op: AFCMP_SUN_S, fmt: OpFormatType_cd_2R},
	AFCOPYSIGN_D:  {mask: 0xffff8000, value: 0x01130000, op: AFCOPYSIGN_D, fmt: OpFormatType_3F},
	AFCOPYSIGN_S:  {mask: 0xffff8000, value: 0x01128000, op: AFCOPYSIGN_S, fmt: OpFormatType_3F},
	AFCVT_D_S:     {mask: 0xfffffc00, value: 0x01192400, op: AFCVT_D_S, fmt: OpFormatType_2F},
	AFCVT_S_D:     {mask: 0xfffffc00, value: 0x01191800, op: AFCVT_S_D, fmt: OpFormatType_2F},
	AFDIV_D:       {mask: 0xffff8000, value: 0x01070000, op: AFDIV_D, fmt: OpFormatType_3F},
	AFDIV_S:       {mask: 0xffff8000, value: 0x01068000, op: AFDIV_S, fmt: OpFormatType_3F},
	AFFINT_D_L:    {mask: 0xfffffc00, value: 0x011d2800, op: AFFINT_D_L, fmt: OpFormatType_2F},
	AFFINT_D_W:    {mask: 0xfffffc00, value: 0x011d2000, op: AFFINT_D_W, fmt: OpFormatType_2F},
	AFFINT_S_L:    {mask: 0xfffffc00, value: 0x011d1800, op: AFFINT_S_L, fmt: OpFormatType_2F},
	AFFINT_S_W:    {mask: 0xfffffc00, value: 0x011d1000, op: AFFINT_S_W, fmt: OpFormatType_2F},
	AFLDGT_D:      {mask: 0xffff8000, value: 0x38748000, op: AFLDGT_D, fmt: OpFormatType_1F_2R},
	AFLDGT_S:      {mask: 0xffff8000, value: 0x38740000, op: AFLDGT_S, fmt: OpFormatType_1F_2R},
	AFLDLE_D:      {mask: 0xffff8000, value: 0x38758000, op: AFLDLE_D, fmt: OpFormatType_1F_2R},
	AFLDLE_S:      {mask: 0xffff8000, value: 0x38750000, op: AFLDLE_S, fmt: OpFormatType_1F_2R},
	AFLDX_D:       {mask: 0xffff8000, value: 0x38340000, op: AFLDX_D, fmt: OpFormatType_1F_2R},
	AFLDX_S:       {mask: 0xffff8000, value: 0x38300000, op: AFLDX_S, fmt: OpFormatType_1F_2R},
	AFLD_D:        {mask: 0xffc00000, value: 0x2b800000, op: AFLD_D, fmt: OpFormatType_2R_si12},
	AFLD_S:        {mask: 0xffc00000, value: 0x2b000000, op: AFLD_S, fmt: OpFormatType_2R_si12},
	AFLOGB_D:      {mask: 0xfffffc00, value: 0x01142800, op: AFLOGB_D, fmt: OpFormatType_2F},
	AFLOGB_S:      {mask: 0xfffffc00, value: 0x01142400, op: AFLOGB_S, fmt: OpFormatType_2F},
	AFMADD_D:      {mask: 0xfff00000, value: 0x08200000, op: AFMADD_D, fmt: OpFormatType_4F},
	AFMADD_S:      {mask: 0xfff00000, value: 0x08100000, op: AFMADD_S, fmt: OpFormatType_4F},
	AFMAXA_D:      {mask: 0xffff8000, value: 0x010d0000, op: AFMAXA_D, fmt: OpFormatType_3F},
	AFMAXA_S:      {mask: 0xffff8000, value: 0x010c8000, op: AFMAXA_S, fmt: OpFormatType_3F},
	AFMAX_D:       {mask: 0xffff8000, value: 0x01090000, op: AFMAX_D, fmt: OpFormatType_3F},
	AFMAX_S:       {mask: 0xffff8000, value: 0x01088000, op: AFMAX_S, fmt: OpFormatType_3F},
	AFMINA_D:      {mask: 0xffff8000, value: 0x010f0000, op: AFMINA_D, fmt: OpFormatType_3F},
	AFMINA_S:      {mask: 0xffff8000, value: 0x010e8000, op: AFMINA_S, fmt: OpFormatType_3F},
	AFMIN_D:       {mask: 0xffff8000, value: 0x010b0000, op: AFMIN_D, fmt: OpFormatType_3F},
	AFMIN_S:       {mask: 0xffff8000, value: 0x010a8000, op: AFMIN_S, fmt: OpFormatType_3F},
	AFMOV_D:       {mask: 0xfffffc00, value: 0x01149800, op: AFMOV_D, fmt: OpFormatType_2F},
	AFMOV_S:       {mask: 0xfffffc00, value: 0x01149400, op: AFMOV_S, fmt: OpFormatType_2F},
	AFMSUB_D:      {mask: 0xfff00000, value: 0x08600000, op: AFMSUB_D, fmt: OpFormatType_4F},
	AFMSUB_S:      {mask: 0xfff00000, value: 0x08500000, op: AFMSUB_S, fmt: OpFormatType_4F},
	AFMUL_D:       {mask: 0xffff8000, value: 0x01050000, op: AFMUL_D, fmt: OpFormatType_3F},
	AFMUL_S:       {mask: 0xffff8000, value: 0x01048000, op: AFMUL_S, fmt: OpFormatType_3F},
	AFNEG_D:       {mask: 0xfffffc00, value: 0x01141800, op: AFNEG_D, fmt: OpFormatType_2F},
	AFNEG_S:       {mask: 0xfffffc00, value: 0x01141400, op: AFNEG_S, fmt: OpFormatType_2F},
	AFNMADD_D:     {mask: 0xfff00000, value: 0x08a00000, op: AFNMADD_D, fmt: OpFormatType_4F},
	AFNMADD_S:     {mask: 0xfff00000, value: 0x08900000, op: AFNMADD_S, fmt: OpFormatType_4F},
	AFNMSUB_D:     {mask: 0xfff00000, value: 0x08e00000, op: AFNMSUB_D, fmt: OpFormatType_4F},
	AFNMSUB_S:     {mask: 0xfff00000, value: 0x08d00000, op: AFNMSUB_S, fmt: OpFormatType_4F},
	AFRECIPE_D:    {mask: 0xfffffc00, value: 0x01147800, op: AFRECIPE_D, fmt: OpFormatType_2F},
	AFRECIPE_S:    {mask: 0xfffffc00, value: 0x01147400, op: AFRECIPE_S, fmt: OpFormatType_2F},
	AFRECIP_D:     {mask: 0xfffffc00, value: 0x01145800, op: AFRECIP_D, fmt: OpFormatType_2F},
	AFRECIP_S:     {mask: 0xfffffc00, value: 0x01145400, op: AFRECIP_S, fmt: OpFormatType_2F},
	AFRINT_D:      {mask: 0xfffffc00, value: 0x011e4800, op: AFRINT_D, fmt: OpFormatType_2F},
	AFRINT_S:      {mask: 0xfffffc00, value: 0x011e4400, op: AFRINT_S, fmt: OpFormatType_2F},
	AFRSQRTE_D:    {mask: 0xfffffc00, value: 0x01148800, op: AFRSQRTE_D, fmt: OpFormatType_2F},
	AFRSQRTE_S:    {mask: 0xfffffc00, value: 0x01148400, op: AFRSQRTE_S, fmt: OpFormatType_2F},
	AFRSQRT_D:     {mask: 0xfffffc00, value: 0x01146800, op: AFRSQRT_D, fmt: OpFormatType_2F},
	AFRSQRT_S:     {mask: 0xfffffc00, value: 0x01146400, op: AFRSQRT_S, fmt: OpFormatType_2F},
	AFSCALEB_D:    {mask: 0xffff8000, value: 0x01110000, op: AFSCALEB_D, fmt: OpFormatType_3F},
	AFSCALEB_S:    {mask: 0xffff8000, value: 0x01108000, op: AFSCALEB_S, fmt: OpFormatType_3F},
	AFSEL:         {mask: 0xfffc0000, value: 0x0d000000, op: AFSEL, fmt: OpFormatType_3R_ca},
	AFSQRT_D:      {mask: 0xfffffc00, value: 0x01144800, op: AFSQRT_D, fmt: OpFormatType_2F},
	AFSQRT_S:      {mask: 0xfffffc00, value: 0x01144400, op: AFSQRT_S, fmt: OpFormatType_2F},
	AFSTGT_D:      {mask: 0xffff8000, value: 0x38768000, op: AFSTGT_D, fmt: OpFormatType_1F_2R},
	AFSTGT_S:      {mask: 0xffff8000, value: 0x38760000, op: AFSTGT_S, fmt: OpFormatType_1F_2R},
	AFSTLE_D:      {mask: 0xffff8000, value: 0x38778000, op: AFSTLE_D, fmt: OpFormatType_1F_2R},
	AFSTLE_S:      {mask: 0xffff8000, value: 0x38770000, op: AFSTLE_S, fmt: OpFormatType_1F_2R},
	AFSTX_D:       {mask: 0xffff8000, value: 0x383c0000, op: AFSTX_D, fmt: OpFormatType_1F_2R},
	AFSTX_S:       {mask: 0xffff8000, value: 0x38380000, op: AFSTX_S, fmt: OpFormatType_1F_2R},
	AFST_D:        {mask: 0xffc00000, value: 0x2bc00000, op: AFST_D, fmt: OpFormatType_2R_si12},
	AFST_S:        {mask: 0xffc00000, value: 0x2b400000, op: AFST_S, fmt: OpFormatType_2R_si12},
	AFSUB_D:       {mask: 0xffff8000, value: 0x01030000, op: AFSUB_D, fmt: OpFormatType_3F},
	AFSUB_S:       {mask: 0xffff8000, value: 0x01028000, op: AFSUB_S, fmt: OpFormatType_3F},
	AFTINTRM_L_D:  {mask: 0xfffffc00, value: 0x011a2800, op: AFTINTRM_L_D, fmt: OpFormatType_2F},
	AFTINTRM_L_S:  {mask: 0xfffffc00, value: 0x011a2400, op: AFTINTRM_L_S, fmt: OpFormatType_2F},
	AFTINTRM_W_D:  {mask: 0xfffffc00, value: 0x011a0800, op: AFTINTRM_W_D, fmt: OpFormatType_2F},
	AFTINTRM_W_S:  {mask: 0xfffffc00, value: 0x011a0400, op: AFTINTRM_W_S, fmt: OpFormatType_2F},
	AFTINTRNE_L_D: {mask: 0xfffffc00, value: 0x011ae800, op: AFTINTRNE_L_D, fmt: OpFormatType_2F},
	AFTINTRNE_L_S: {mask: 0xfffffc00, value: 0x011ae400, op: AFTINTRNE_L_S, fmt: OpFormatType_2F},
	AFTINTRNE_W_D: {mask: 0xfffffc00, value: 0x011ac800, op: AFTINTRNE_W_D, fmt: OpFormatType_2F},
	AFTINTRNE_W_S: {mask: 0xfffffc00, value: 0x011ac400, op: AFTINTRNE_W_S, fmt: OpFormatType_2F},
	AFTINTRP_L_D:  {mask: 0xfffffc00, value: 0x011a6800, op: AFTINTRP_L_D, fmt: OpFormatType_2F},
	AFTINTRP_L_S:  {mask: 0xfffffc00, value: 0x011a6400, op: AFTINTRP_L_S, fmt: OpFormatType_2F},
	AFTINTRP_W_D:  {mask: 0xfffffc00, value: 0x011a4800, op: AFTINTRP_W_D, fmt: OpFormatType_2F},
	AFTINTRP_W_S:  {mask: 0xfffffc00, value: 0x011a4400, op: AFTINTRP_W_S, fmt: OpFormatType_2F},
	AFTINTRZ_L_D:  {mask: 0xfffffc00, value: 0x011aa800, op: AFTINTRZ_L_D, fmt: OpFormatType_2F},
	AFTINTRZ_L_S:  {mask: 0xfffffc00, value: 0x011aa400, op: AFTINTRZ_L_S, fmt: OpFormatType_2F},
	AFTINTRZ_W_D:  {mask: 0xfffffc00, value: 0x011a8800, op: AFTINTRZ_W_D, fmt: OpFormatType_2F},
	AFTINTRZ_W_S:  {mask: 0xfffffc00, value: 0x011a8400, op: AFTINTRZ_W_S, fmt: OpFormatType_2F},
	AFTINT_L_D:    {mask: 0xfffffc00, value: 0x011b2800, op: AFTINT_L_D, fmt: OpFormatType_2F},
	AFTINT_L_S:    {mask: 0xfffffc00, value: 0x011b2400, op: AFTINT_L_S, fmt: OpFormatType_2F},
	AFTINT_W_D:    {mask: 0xfffffc00, value: 0x011b0800, op: AFTINT_W_D, fmt: OpFormatType_2F},
	AFTINT_W_S:    {mask: 0xfffffc00, value: 0x011b0400, op: AFTINT_W_S, fmt: OpFormatType_2F},
	AIBAR:         {mask: 0xffff8000, value: 0x38728000, op: AIBAR, fmt: OpFormatType_hint},
	AIDLE:         {mask: 0xffff8000, value: 0x06488000, op: AIDLE, fmt: OpFormatType_level},
	AINVTLB:       {mask: 0xffff8000, value: 0x06498000, op: AINVTLB, fmt: OpFormatType_op_2R},
	AIOCSRRD_B:    {mask: 0xfffffc00, value: 0x06480000, op: AIOCSRRD_B, fmt: OpFormatType_2R},
	AIOCSRRD_D:    {mask: 0xfffffc00, value: 0x06480c00, op: AIOCSRRD_D, fmt: OpFormatType_2R},
	AIOCSRRD_H:    {mask: 0xfffffc00, value: 0x06480400, op: AIOCSRRD_H, fmt: OpFormatType_2R},
	AIOCSRRD_W:    {mask: 0xfffffc00, value: 0x06480800, op: AIOCSRRD_W, fmt: OpFormatType_2R},
	AIOCSRWR_B:    {mask: 0xfffffc00, value: 0x06481000, op: AIOCSRWR_B, fmt: OpFormatType_2R},
	AIOCSRWR_D:    {mask: 0xfffffc00, value: 0x06481c00, op: AIOCSRWR_D, fmt: OpFormatType_2R},
	AIOCSRWR_H:    {mask: 0xfffffc00, value: 0x06481400, op: AIOCSRWR_H, fmt: OpFormatType_2R},
	AIOCSRWR_W:    {mask: 0xfffffc00, value: 0x06481800, op: AIOCSRWR_W, fmt: OpFormatType_2R},
	AJIRL:         {mask: 0xfc000000, value: 0x4c000000, op: AJIRL, fmt: OpFormatType_rd_rj_offset},
	ALDDIR:        {mask: 0xfffc0000, value: 0x06400000, op: ALDDIR, fmt: OpFormatType_2R_level},
	ALDGT_B:       {mask: 0xffff8000, value: 0x38780000, op: ALDGT_B, fmt: OpFormatType_3R},
	ALDGT_D:       {mask: 0xffff8000, value: 0x38798000, op: ALDGT_D, fmt: OpFormatType_3R},
	ALDGT_H:       {mask: 0xffff8000, value: 0x38788000, op: ALDGT_H, fmt: OpFormatType_3R},
	ALDGT_W:       {mask: 0xffff8000, value: 0x38790000, op: ALDGT_W, fmt: OpFormatType_3R},
	ALDLE_B:       {mask: 0xffff8000, value: 0x387a0000, op: ALDLE_B, fmt: OpFormatType_3R},
	ALDLE_D:       {mask: 0xffff8000, value: 0x387b8000, op: ALDLE_D, fmt: OpFormatType_3R},
	ALDLE_H:       {mask: 0xffff8000, value: 0x387a8000, op: ALDLE_H, fmt: OpFormatType_3R},
	ALDLE_W:       {mask: 0xffff8000, value: 0x387b0000, op: ALDLE_W, fmt: OpFormatType_3R},
	ALDPTE:        {mask: 0xfffc001f, value: 0x06440000, op: ALDPTE, fmt: OpFormatType_0_1R_seq},
	ALDPTR_D:      {mask: 0xff000000, value: 0x26000000, op: ALDPTR_D, fmt: OpFormatType_2R_si14},
	ALDPTR_W:      {mask: 0xff000000, value: 0x24000000, op: ALDPTR_W, fmt: OpFormatType_2R_si14},
	ALDX_B:        {mask: 0xffff8000, value: 0x38000000, op: ALDX_B, fmt: OpFormatType_3R},
	ALDX_BU:       {mask: 0xffff8000, value: 0x38200000, op: ALDX_BU, fmt: OpFormatType_3R},
	ALDX_D:        {mask: 0xffff8000, value: 0x380c0000, op: ALDX_D, fmt: OpFormatType_3R},
	ALDX_H:        {mask: 0xffff8000, value: 0x38040000, op: ALDX_H, fmt: OpFormatType_3R},
	ALDX_HU:       {mask: 0xffff8000, value: 0x38240000, op: ALDX_HU, fmt: OpFormatType_3R},
	ALDX_W:        {mask: 0xffff8000, value: 0x38080000, op: ALDX_W, fmt: OpFormatType_3R},
	ALDX_WU:       {mask: 0xffff8000, value: 0x38280000, op: ALDX_WU, fmt: OpFormatType_3R},
	ALD_B:         {mask: 0xffc00000, value: 0x28000000, op: ALD_B, fmt: OpFormatType_2R_si12},
	ALD_BU:        {mask: 0xffc00000, value: 0x2a000000, op: ALD_BU, fmt: OpFormatType_2R_si12},
	ALD_D:         {mask: 0xffc00000, value: 0x28c00000, op: ALD_D, fmt: OpFormatType_2R_si12},
	ALD_H:         {mask: 0xffc00000, value: 0x28400000, op: ALD_H, fmt: OpFormatType_2R_si12},
	ALD_HU:        {mask: 0xffc00000, value: 0x2a400000, op: ALD_HU, fmt: OpFormatType_2R_si12},
	ALD_W:         {mask: 0xffc00000, value: 0x28800000, op: ALD_W, fmt: OpFormatType_2R_si12},
	ALD_WU:        {mask: 0xffc00000, value: 0x2a800000, op: ALD_WU, fmt: OpFormatType_2R_si12},
	ALLACQ_D:      {mask: 0xfffffc00, value: 0x38578800, op: ALLACQ_D, fmt: OpFormatType_2R},
	ALLACQ_W:      {mask: 0xfffffc00, value: 0x38578000, op: ALLACQ_W, fmt: OpFormatType_2R},
	ALL_D:         {mask: 0xff000000, value: 0x22000000, op: ALL_D, fmt: OpFormatType_2R_si14},
	ALL_W:         {mask: 0xff000000, value: 0x20000000, op: ALL_W, fmt: OpFormatType_2R_si14},
	ALU12I_W:      {mask: 0xfe000000, value: 0x14000000, op: ALU12I_W, fmt: OpFormatType_1R_si20},
	ALU32I_D:      {mask: 0xfe000000, value: 0x16000000, op: ALU32I_D, fmt: OpFormatType_1R_si20},
	ALU52I_D:      {mask: 0xffc00000, value: 0x03000000, op: ALU52I_D, fmt: OpFormatType_2R_si12},
	AMASKEQZ:      {mask: 0xffff8000, value: 0x00130000, op: AMASKEQZ, fmt: OpFormatType_3R},
	AMASKNEZ:      {mask: 0xffff8000, value: 0x00138000, op: AMASKNEZ, fmt: OpFormatType_3R},
	AMOD_D:        {mask: 0xffff8000, value: 0x00228000, op: AMOD_D, fmt: OpFormatType_3R},
	AMOD_DU:       {mask: 0xffff8000, value: 0x00238000, op: AMOD_DU, fmt: OpFormatType_3R},
	AMOD_W:        {mask: 0xffff8000, value: 0x00208000, op: AMOD_W, fmt: OpFormatType_3R},
	AMOD_WU:       {mask: 0xffff8000, value: 0x00218000, op: AMOD_WU, fmt: OpFormatType_3R},
	AMOVCF2FR:     {mask: 0xffffff00, value: 0x0114d400, op: AMOVCF2FR, fmt: OpFormatType_1R_cj},
	AMOVCF2GR:     {mask: 0xffffff00, value: 0x0114dc00, op: AMOVCF2GR, fmt: OpFormatType_1R_cj},
	AMOVFCSR2GR:   {mask: 0xfffffc00, value: 0x0114c800, op: AMOVFCSR2GR, fmt: OpFormatType_1R_fcsr},
	AMOVFR2CF:     {mask: 0xfffffc18, value: 0x0114d000, op: AMOVFR2CF, fmt: OpFormatType_cd_1R},
	AMOVFR2GR_D:   {mask: 0xfffffc00, value: 0x0114b800, op: AMOVFR2GR_D, fmt: OpFormatType_1R_1F},
	AMOVFR2GR_S:   {mask: 0xfffffc00, value: 0x0114b400, op: AMOVFR2GR_S, fmt: OpFormatType_1R_1F},
	AMOVFRH2GR_S:  {mask: 0xfffffc00, value: 0x0114bc00, op: AMOVFRH2GR_S, fmt: OpFormatType_1R_1F},
	AMOVGR2CF:     {mask: 0xfffffc18, value: 0x0114d800, op: AMOVGR2CF, fmt: OpFormatType_cd_1R},
	AMOVGR2FCSR:   {mask: 0xfffffc00, value: 0x0114c000, op: AMOVGR2FCSR, fmt: OpFormatType_fcsr_1R},
	AMOVGR2FRH_W:  {mask: 0xfffffc00, value: 0x0114ac00, op: AMOVGR2FRH_W, fmt: OpFormatType_1F_1R},
	AMOVGR2FR_D:   {mask: 0xfffffc00, value: 0x0114a800, op: AMOVGR2FR_D, fmt: OpFormatType_1F_1R},
	AMOVGR2FR_W:   {mask: 0xfffffc00, value: 0x0114a400, op: AMOVGR2FR_W, fmt: OpFormatType_1F_1R},
	AMULH_D:       {mask: 0xffff8000, value: 0x001e0000, op: AMULH_D, fmt: OpFormatType_3R},
	AMULH_DU:      {mask: 0xffff8000, value: 0x001e8000, op: AMULH_DU, fmt: OpFormatType_3R},
	AMULH_W:       {mask: 0xffff8000, value: 0x001c8000, op: AMULH_W, fmt: OpFormatType_3R},
	AMULH_WU:      {mask: 0xffff8000, value: 0x001d0000, op: AMULH_WU, fmt: OpFormatType_3R},
	AMULW_D_W:     {mask: 0xffff8000, value: 0x001f0000, op: AMULW_D_W, fmt: OpFormatType_3R},
	AMULW_D_WU:    {mask: 0xffff8000, value: 0x001f8000, op: AMULW_D_WU, fmt: OpFormatType_3R},
	AMUL_D:        {mask: 0xffff8000, value: 0x001d8000, op: AMUL_D, fmt: OpFormatType_3R},
	AMUL_W:        {mask: 0xffff8000, value: 0x001c0000, op: AMUL_W, fmt: OpFormatType_3R},
	ANOR:          {mask: 0xffff8000, value: 0x00140000, op: ANOR, fmt: OpFormatType_3R},
	AOR:           {mask: 0xffff8000, value: 0x00150000, op: AOR, fmt: OpFormatType_3R},
	AORI:          {mask: 0xffc00000, value: 0x03800000, op: AORI, fmt: OpFormatType_2R_ui12},
	AORN:          {mask: 0xffff8000, value: 0x00160000, op: AORN, fmt: OpFormatType_3R},
	APCADDI:       {mask: 0xfe000000, value: 0x18000000, op: APCADDI, fmt: OpFormatType_1R_si20},
	APCADDU12I:    {mask: 0xfe000000, value: 0x1c000000, op: APCADDU12I, fmt: OpFormatType_1R_si20},
	APCADDU18I:    {mask: 0xfe000000, value: 0x1e000000, op: APCADDU18I, fmt: OpFormatType_1R_si20},
	APCALAU12I:    {mask: 0xfe000000, value: 0x1a000000, op: APCALAU12I, fmt: OpFormatType_1R_si20},
	APRELD:        {mask: 0xffc00000, value: 0x2ac00000, op: APRELD, fmt: OpFormatType_hint_1R_si12},
	APRELDX:       {mask: 0xffff8000, value: 0x382c0000, op: APRELDX, fmt: OpFormatType_hint_2R},
	ARDTIMEH_W:    {mask: 0xfffffc00, value: 0x00006400, op: ARDTIMEH_W, fmt: OpFormatType_2R},
	ARDTIMEL_W:    {mask: 0xfffffc00, value: 0x00006000, op: ARDTIMEL_W, fmt: OpFormatType_2R},
	ARDTIME_D:     {mask: 0xfffffc00, value: 0x00006800, op: ARDTIME_D, fmt: OpFormatType_2R},
	AREVB_2H:      {mask: 0xfffffc00, value: 0x00003000, op: AREVB_2H, fmt: OpFormatType_2R},
	AREVB_2W:      {mask: 0xfffffc00, value: 0x00003800, op: AREVB_2W, fmt: OpFormatType_2R},
	AREVB_4H:      {mask: 0xfffffc00, value: 0x00003400, op: AREVB_4H, fmt: OpFormatType_2R},
	AREVB_D:       {mask: 0xfffffc00, value: 0x00003c00, op: AREVB_D, fmt: OpFormatType_2R},
	AREVH_2W:      {mask: 0xfffffc00, value: 0x00004000, op: AREVH_2W, fmt: OpFormatType_2R},
	AREVH_D:       {mask: 0xfffffc00, value: 0x00004400, op: AREVH_D, fmt: OpFormatType_2R},
	AROTRI_D:      {mask: 0xffff0000, value: 0x004d0000, op: AROTRI_D, fmt: OpFormatType_2R_ui6},
	AROTRI_W:      {mask: 0xffff8000, value: 0x004c8000, op: AROTRI_W, fmt: OpFormatType_2R_ui5},
	AROTR_D:       {mask: 0xffff8000, value: 0x001b8000, op: AROTR_D, fmt: OpFormatType_3R},
	AROTR_W:       {mask: 0xffff8000, value: 0x001b0000, op: AROTR_W, fmt: OpFormatType_3R},
	ASCREL_D:      {mask: 0xfffffc00, value: 0x38578c00, op: ASCREL_D, fmt: OpFormatType_2R},
	ASCREL_W:      {mask: 0xfffffc00, value: 0x38578400, op: ASCREL_W, fmt: OpFormatType_2R},
	ASC_D:         {mask: 0xff000000, value: 0x23000000, op: ASC_D, fmt: OpFormatType_2R_si14},
	ASC_Q:         {mask: 0xffff8000, value: 0x38570000, op: ASC_Q, fmt: OpFormatType_3R},
	ASC_W:         {mask: 0xff000000, value: 0x21000000, op: ASC_W, fmt: OpFormatType_2R_si14},
	ASLLI_D:       {mask: 0xffff0000, value: 0x00410000, op: ASLLI_D, fmt: OpFormatType_2R_ui6},
	ASLLI_W:       {mask: 0xffff8000, value: 0x00408000, op: ASLLI_W, fmt: OpFormatType_2R_ui5},
	ASLL_D:        {mask: 0xffff8000, value: 0x00188000, op: ASLL_D, fmt: OpFormatType_3R},
	ASLL_W:        {mask: 0xffff8000, value: 0x00170000, op: ASLL_W, fmt: OpFormatType_3R},
	ASLT:          {mask: 0xffff8000, value: 0x00120000, op: ASLT, fmt: OpFormatType_3R},
	ASLTI:         {mask: 0xffc00000, value: 0x02000000, op: ASLTI, fmt: OpFormatType_2R_si12},
	ASLTU:         {mask: 0xffff8000, value: 0x00128000, op: ASLTU, fmt: OpFormatType_3R},
	ASLTUI:        {mask: 0xffc00000, value: 0x02400000, op: ASLTUI, fmt: OpFormatType_2R_si12},
	ASRAI_D:       {mask: 0xffff0000, value: 0x00490000, op: ASRAI_D, fmt: OpFormatType_2R_ui6},
	ASRAI_W:       {mask: 0xffff8000, value: 0x00488000, op: ASRAI_W, fmt: OpFormatType_2R_ui5},
	ASRA_D:        {mask: 0xffff8000, value: 0x00198000, op: ASRA_D, fmt: OpFormatType_3R},
	ASRA_W:        {mask: 0xffff8000, value: 0x00180000, op: ASRA_W, fmt: OpFormatType_3R},
	ASRLI_D:       {mask: 0xffff0000, value: 0x00450000, op: ASRLI_D, fmt: OpFormatType_2R_ui6},
	ASRLI_W:       {mask: 0xffff8000, value: 0x00448000, op: ASRLI_W, fmt: OpFormatType_2R_ui5},
	ASRL_D:        {mask: 0xffff8000, value: 0x00190000, op: ASRL_D, fmt: OpFormatType_3R},
	ASRL_W:        {mask: 0xffff8000, value: 0x00178000, op: ASRL_W, fmt: OpFormatType_3R},
	ASTGT_B:       {mask: 0xffff8000, value: 0x387c0000, op: ASTGT_B, fmt: OpFormatType_3R},
	ASTGT_D:       {mask: 0xffff8000, value: 0x387d8000, op: ASTGT_D, fmt: OpFormatType_3R},
	ASTGT_H:       {mask: 0xffff8000, value: 0x387c8000, op: ASTGT_H, fmt: OpFormatType_3R},
	ASTGT_W:       {mask: 0xffff8000, value: 0x387d0000, op: ASTGT_W, fmt: OpFormatType_3R},
	ASTLE_B:       {mask: 0xffff8000, value: 0x387e0000, op: ASTLE_B, fmt: OpFormatType_3R},
	ASTLE_D:       {mask: 0xffff8000, value: 0x387f8000, op: ASTLE_D, fmt: OpFormatType_3R},
	ASTLE_H:       {mask: 0xffff8000, value: 0x387e8000, op: ASTLE_H, fmt: OpFormatType_3R},
	ASTLE_W:       {mask: 0xffff8000, value: 0x387f0000, op: ASTLE_W, fmt: OpFormatType_3R},
	ASTPTR_D:      {mask: 0xff000000, value: 0x27000000, op: ASTPTR_D, fmt: OpFormatType_2R_si14},
	ASTPTR_W:      {mask: 0xff000000, value: 0x25000000, op: ASTPTR_W, fmt: OpFormatType_2R_si14},
	ASTX_B:        {mask: 0xffff8000, value: 0x38100000, op: ASTX_B, fmt: OpFormatType_3R},
	ASTX_D:        {mask: 0xffff8000, value: 0x381c0000, op: ASTX_D, fmt: OpFormatType_3R},
	ASTX_H:        {mask: 0xffff8000, value: 0x38140000, op: ASTX_H, fmt: OpFormatType_3R},
	ASTX_W:        {mask: 0xffff8000, value: 0x38180000, op: ASTX_W, fmt: OpFormatType_3R},
	AST_B:         {mask: 0xffc00000, value: 0x29000000, op: AST_B, fmt: OpFormatType_2R_si12},
	AST_D:         {mask: 0xffc00000, value: 0x29c00000, op: AST_D, fmt: OpFormatType_2R_si12},
	AST_H:         {mask: 0xffc00000, value: 0x29400000, op: AST_H, fmt: OpFormatType_2R_si12},
	AST_W:         {mask: 0xffc00000, value: 0x29800000, op: AST_W, fmt: OpFormatType_2R_si12},
	ASUB_D:        {mask: 0xffff8000, value: 0x00118000, op: ASUB_D, fmt: OpFormatType_3R},
	ASUB_W:        {mask: 0xffff8000, value: 0x00110000, op: ASUB_W, fmt: OpFormatType_3R},
	ASYSCALL:      {mask: 0xffff8000, value: 0x002b0000, op: ASYSCALL, fmt: OpFormatType_code},
	ATLBCLR:       {mask: 0xffffffff, value: 0x06482000, op: ATLBCLR, fmt: OpFormatType_NULL},
	ATLBFILL:      {mask: 0xffffffff, value: 0x06483400, op: ATLBFILL, fmt: OpFormatType_NULL},
	ATLBFLUSH:     {mask: 0xffffffff, value: 0x06482400, op: ATLBFLUSH, fmt: OpFormatType_NULL},
	ATLBRD:        {mask: 0xffffffff, value: 0x06482c00, op: ATLBRD, fmt: OpFormatType_NULL},
	ATLBSRCH:      {mask: 0xffffffff, value: 0x06482800, op: ATLBSRCH, fmt: OpFormatType_NULL},
	ATLBWR:        {mask: 0xffffffff, value: 0x06483000, op: ATLBWR, fmt: OpFormatType_NULL},
	AXOR:          {mask: 0xffff8000, value: 0x00158000, op: AXOR, fmt: OpFormatType_3R},
	AXORI:         {mask: 0xffc00000, value: 0x03c00000, op: AXORI, fmt: OpFormatType_2R_ui12},
}
