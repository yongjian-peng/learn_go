package iterator

// 班级迭代器

type Iterator interface {
	Next() Member  // 迭代器下一个成员
	HasMore() bool // 是否还有
}

// Class 班级，包括老师和学生
type Class struct {
	name     string
	teacher  *Teacher
	students []*Student
}

// memberIterator 班级成员迭代器实现
type memberIterator struct {
	class *Class // 需要迭代的班级
	index int    // 迭代索引
}

func (m *memberIterator) Next() Member {
	// 迭代索引为 -1 时，返回老师成员，否则遍历学生slice
	if m.index == -1 {
		m.index++
		return m.class.teacher
	}
	student := m.class.students[m.index]
	m.index++
	return student
}

func (m *memberIterator) HasMore() bool {
	return m.index < len(m.class.students)
}

// Iterable 可迭代集合接口，实现此接口返回迭代器
type Iterable interface {
	CreateIterator() Iterator
}

// NewClass 根据班主任老师名称，授课创建班级
func NewClass(name, teacherName, teacherSubject string) *Class {
	return &Class{
		name:    name,
		teacher: NewTeacher(teacherName, teacherSubject),
	}
}

// CreateIterator 创建班级迭代器
func (c *Class) CreateIterator() Iterator {
	return &memberIterator{
		class: c,
		index: -1, // 迭代索引初始化为-1， 从老师开始迭代
	}
}

func (c *Class) Name() string {
	return c.name
}

// AddStudent 班级添加同学
func (c *Class) AddStudent(students ...*Student) {
	c.students = append(c.students, students...)
}
