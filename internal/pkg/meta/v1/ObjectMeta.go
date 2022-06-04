package v1

type ObjectMeta struct {
	ID uint64 `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENt;column:id;comment:主键"`
}
