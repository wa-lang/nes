
/*
创建一个空 JS 对象，其将对象存入对象池后，返回ID（既对象句柄）
*/
#wa:import extobj new_obj
func jsNewObj() => u32

/*
释放对象池中的指定对象，如果该对象持有其它资源（比如GpuTexture），应执行相应清理动作
当 h 指定的对象不再被凹内部引用时，该方法将被触发
当对象被释放后，其在对象池中的 id 应被回收，后续创建新对象时可复用，避免对象池无限增大
*/
#wa:import extobj free_obj
func jsFreeObj(h: u32)

/*
向一个对象中添加 i32 型的成员
 h: 对象在对象池中的id
 member_name: 成员名字
 value: 成员的值
*/
#wa:import extobj set_member_i32
func jsSetMember_i32(h: u32, member_name: string, value: i32)

/*
向一个对象中添加 string 型的成员
 h: 对象在对象池中的id
 member_name: 成员名字
 value: 成员的值
*/
#wa:import extobj set_member_string
func jsSetMember_string(h: u32, member_name: string, value: string)

/*
将源对象添加为目标对象的成员
 h: 目标对象在对象池中的id
 member_name: 成员名字
 value: 源对象在对象池中的id
*/
#wa:import extobj set_member_obj
func jsSetMember_obj(h: u32, member_name: string, value: u32)

/*
创建一个空 JS 数组，其将数组存入对象池后，返回ID（既对象句柄）
*/
#wa:import extobj new_array
func jsNewArray() => u32

/*
向数组中添加 i32 型元素
 h: 数组对象在对象池中的id
 value: 被添加的值
*/
#wa:import extobj append_i32
func jsAppend_i32(h: u32, value: i32)

/*
向数组中添加 string 型元素
 h: 数组对象在对象池中的id
 value: 被添加的值
*/
#wa:import extobj append_string
func jsAppend_string(h: u32, value: string)

/*
向目标数组中添加源对象元素
 h: 目标数组对象在对象池中的id
 value: 源对象在对象池中的id
*/
#wa:import extobj append_obj
func jsAppend_obj(h: u32, value: u32)

