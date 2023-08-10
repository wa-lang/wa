
(func $$wa.runtime.queryIface (param $d.b i32) (param $d.d i32) (param $itab i32) (param $eq i32) (param $ihash i32)
  (result i32 i32 i32 i32)
	(local $t i32)
	local.get $itab
	if (result i32 i32 i32 i32)
	  local.get $itab
	  i32.load offset=0 align=4
	  local.get $ihash
	  i32.const 0
	  call $$wa.runtime.getItab
	  local.set $t
	  local.get $t
	  if (result i32 i32 i32 i32)
	    local.get $d.b
	    local.get $d.d
	    local.get $t
	    local.get $eq
	  else
	    local.get $d.b
	    call $$Release
	    i32.const 0
	    i32.const 0
	    i32.const 0
	    i32.const 0
	    unreachable
	  end
	else
	  local.get $d.b
	  call $$Release
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
	  call $$wa.runtime.getItab
	  local.set $t
	  local.get $t
	  if (result i32 i32 i32 i32 i32)
	    local.get $d.b
	    local.get $d.d
	    local.get $t
	    local.get $eq
	    i32.const 1
	  else
	    local.get $d.b
	    call $$Release
	    i32.const 0
	    i32.const 0
	    i32.const 0
	    i32.const 0
	    i32.const 0
	  end
	else
	  local.get $d.b
	  call $$Release
	  i32.const 0
	  i32.const 0
	  i32.const 0
	  i32.const 0
	  i32.const 0
	end
)
