package x86

import (
	"log"

	"wa-lang.org/wa/internal/p9asm/obj"
)

func obj_Addrel(s *obj.LSym) *obj.Reloc {
	s.R = append(s.R, obj.Reloc{})
	return &s.R[len(s.R)-1]
}

// 增长符号对应的机器码列表
func obj_Symgrow(s *obj.LSym, lsiz int64) {
	siz := int(lsiz)
	if int64(siz) != lsiz {
		log.Fatalf("Symgrow size %d too long", lsiz)
	}
	if len(s.P) >= siz {
		return
	}
	for cap(s.P) < siz {
		s.P = append(s.P[:cap(s.P)], 0)
	}
	s.P = s.P[:siz]
}

func obj_Appendp(ctxt *obj.Link, q *obj.Prog) *obj.Prog {
	p := obj.NewProg(ctxt)
	p.Link = q.Link
	q.Link = p
	p.Lineno = q.Lineno
	p.Mode = q.Mode
	return p
}

func obj_Brchain(ctxt *obj.Link, p *obj.Prog) *obj.Prog {
	for i := 0; i < 20; i++ {
		if p == nil || p.As != obj.AJMP || p.Pcond == nil {
			return p
		}
		p = p.Pcond
	}

	return nil
}
