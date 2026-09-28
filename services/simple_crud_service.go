package services

import (
	"fmt"
	"reflect"

	"github.com/beego/beego/v2/client/orm"
)

type SimpleCRUDService[T any] struct {
	FilterMap map[string]string
}

func (s SimpleCRUDService[T]) Create(input *T) (*T, error) {
	setBoolField(input, "Activo", true)
	o := orm.NewOrm()
	if _, err := o.Insert(input); err != nil {
		return nil, fmt.Errorf("insert %T: %w", input, err)
	}
	id := int64Field(input, "Id")
	if id == 0 {
		return input, nil
	}
	return s.GetByID(id)
}

func (s SimpleCRUDService[T]) GetByID(id int64) (*T, error) {
	o := orm.NewOrm()
	item := new(T)
	if err := o.QueryTable(item).Filter("id", id).RelatedSel().One(item); err != nil {
		return nil, fmt.Errorf("read %T %d: %w", item, id, err)
	}
	return item, nil
}

func (s SimpleCRUDService[T]) List(rawQuery string) ([]*T, error) {
	filters, err := ParseQueryFilters(rawQuery)
	if err != nil {
		return nil, err
	}

	model := new(T)
	qs := orm.NewOrm().QueryTable(model)
	for key, value := range filters {
		ormField := key
		if s.FilterMap != nil {
			mapped, ok := s.FilterMap[key]
			if !ok {
				return nil, fmt.Errorf("unsupported filter %q", key)
			}
			ormField = mapped
		}
		qs = qs.Filter(ormField, value)
	}

	var items []*T
	if _, err := qs.RelatedSel().All(&items); err != nil {
		return nil, fmt.Errorf("list %T: %w", model, err)
	}
	return items, nil
}

func (s SimpleCRUDService[T]) Update(id int64, input *T) (*T, error) {
	setInt64Field(input, "Id", id)
	o := orm.NewOrm()
	if _, err := o.Update(input); err != nil {
		return nil, fmt.Errorf("update %T %d: %w", input, id, err)
	}
	return s.GetByID(id)
}

func (s SimpleCRUDService[T]) SoftDelete(id int64) error {
	current, err := s.GetByID(id)
	if err != nil {
		return err
	}
	setBoolField(current, "Activo", false)
	if _, err := orm.NewOrm().Update(current, "activo", "fecha_modificacion"); err != nil {
		return fmt.Errorf("soft delete %T %d: %w", current, id, err)
	}
	return nil
}

func setInt64Field[T any](input *T, name string, value int64) {
	field := reflect.ValueOf(input).Elem().FieldByName(name)
	if field.IsValid() && field.CanSet() && field.Kind() == reflect.Int64 {
		field.SetInt(value)
	}
}

func int64Field[T any](input *T, name string) int64 {
	field := reflect.ValueOf(input).Elem().FieldByName(name)
	if field.IsValid() && field.Kind() == reflect.Int64 {
		return field.Int()
	}
	return 0
}

func setBoolField[T any](input *T, name string, value bool) {
	field := reflect.ValueOf(input).Elem().FieldByName(name)
	if field.IsValid() && field.CanSet() && field.Kind() == reflect.Bool {
		field.SetBool(value)
	}
}
