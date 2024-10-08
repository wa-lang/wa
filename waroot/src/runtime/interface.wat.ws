
(func $$wa.runtime.queryIface (param $d.b i32) (param $d.d i32) (param $itab i32) (param $eq i32) (param $ihash i32)
  (result i32 i32 i32 i32)
	(local $t i32)
	local.get $itab
	if (result i32 i32 i32 i32)
	  local.get $itab
	  i32.load offset=0 align=4
	  local.get $ihash
	  i32.const 0
	  call $runtime.getItab
	  local.set $t
	  local.get $t
	  if (result i32 i32 i32 i32)
	    local.get $d.b
		call $runtime.Block.Retain
	    local.get $d.d
	    local.get $t
	    local.get $eq
	  else
	    i32.const 0
	    i32.const 0
	    i32.const 0
	    i32.const 0
	    unreachable
	  end
	else
	  i32.const 0
	  i32.const 0
	  i32.const 0
	  i32.const 0
	  unreachable
	end
)

(func $$wa.runtime.queryIface_CommaOk (param $d.b i32) (param $d.d i32) (param $itab i32) (param $eq i32) (param $ihash i32)
  (result i32 i32 i32 i32 i32)
	(local $t i32)
	local.get $itab
	if (result i32 i32 i32 i32 i32)
	  local.get $itab
	  i32.load offset=0 align=4
	  local.get $ihash
	  i32.const 1
	  call $runtime.getItab
	  local.set $t
	  local.get $t
	  if (result i32 i32 i32 i32 i32)
	    local.get $d.b
		call $runtime.Block.Retain
	    local.get $d.d
	    local.get $t
	    local.get $eq
	    i32.const 1
	  else
	    i32.const 0
	    i32.const 0
	    i32.const 0
	    i32.const 0
	    i32.const 0
	  end
	else
	  i32.const 0
	  i32.const 0
	  i32.const 0
	  i32.const 0
	  i32.const 0
	end
)

(func $runtime.Compare (param $l.d.b i32) (param $l.d.d i32) (param $l.itab i32) (param $l.comp i32) (param $r.d.b i32) (param $r.d.d i32) (param $r.itab i32) (param $r.comp i32)
  (result i32)
  local.get $l.comp
  local.get $r.comp
  i32.lt_s
  if (result i32) ;;if l.comp < r.comp
    i32.const -1
  else
    local.get $l.comp
	local.get $r.comp
	i32.gt_s
	if (result i32) ;;if l.comp > r.comp
	  i32.const 1
	else ;;if l.comp == r.comp:
	  local.get $l.comp
	  if (result i32) ;;if comp != 0, compare by type.comp:
	    local.get $l.d.d
	    local.get $r.d.d
	    local.get $l.comp
		call_indirect (type $$wa.runtime.comp)
	  else ;;if comp == 0, compare as ref:
	    local.get $l.d.d
		local.get $r.d.d
		i32.lt_u
		if (result i32)
		  i32.const -1
		else
		  local.get $l.d.d
		  local.get $r.d.d
		  i32.gt_u
		end
	  end
	end
  end
)
