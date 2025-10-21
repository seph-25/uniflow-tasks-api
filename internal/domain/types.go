package domain

// Task Status
const (
    StatusTodo       = "todo"
    StatusInProgress = "in-progress"
    StatusInReview   = "in-review"
    StatusDone       = "done"
    StatusCancelled  = "cancelled"
)

var ValidStatuses = []string{StatusTodo, StatusInProgress, StatusInReview, StatusDone, StatusCancelled}

// Task Priority
const (
    PriorityLow    = "low"
    PriorityMedium = "medium"
    PriorityHigh   = "high"
    PriorityUrgent = "urgent"
)

var ValidPriorities = []string{PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent}

// Task Type
const (
    TypeAssignment = "assignment"
    TypeExam       = "exam"
    TypeReading    = "reading"
    TypePresentation = "presentation"
    TypeLab        = "lab"
    TypeQuiz       = "quiz"
    TypeEssay      = "essay"
    TypeGroupWork  = "group-work"
)

var ValidTypes = []string{TypeAssignment, TypeExam, TypeReading, TypePresentation, TypeLab, TypeQuiz, TypeEssay, TypeGroupWork}