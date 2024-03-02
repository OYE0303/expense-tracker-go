package efactory

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

// Config is the configuration for the factory
type Config[T any] struct {
	// BluePrint is a client-defined function to create a new value
	// if not provided, the factory will set non-zero values to the value
	// BluePrint must follow the below signature
	// type bluePrintFunc[T any] func(i int, last T) T
	BluePrint bluePrintFunc[T]

	// Inserter is a client-defined function to insert a value into the database
	// if not provided, the factory will use the default inserter
	// Inserter must follow the below signature
	// type inserter[T any] func(db *sql.DB, v T) (T, error)
	Inserter inserter[T]

	// DB is the database connection
	// must be provided if using Insert or InsertWithAss
	DB *sql.DB

	// TableName is the table name for the value
	// must be provided if not providing Inserter
	TableName string
}

type Factory[T any] struct {
	db        *sql.DB
	bluePrint bluePrintFunc[T]
	inserter  inserter[T]
	tableName string
	dataType  reflect.Type
	empty     T
	index     int

	// map from name to trait function
	traits map[string]setTraiter[T]

	// map from name to list of associations
	// e.g. "User" -> []*User
	associations map[string][]interface{}

	// map from tag to metadata
	// e.g. "User" -> {tableName: "users", fieldName: "UserID"}
	tagToInfo map[string]tagInfo
}

// bluePrintFunc is a client-defined function to create a new value
type bluePrintFunc[T any] func(i int, last T) T

// inserter is a client-defined function to insert a value into the database
type inserter[T any] func(db *sql.DB, v T) (T, error)

// SetTrait is a client-defined function to add a trait to mutate the value
type setTraiter[T any] func(v *T)

// tagInfo is the metadata for the tag
type tagInfo struct {
	tableName string
	fieldName string
}

// builder is for building a single value
type builder[T any] struct {
	v      *T
	errors []error
	f      *Factory[T]
}

// builderList is for building a list of values
type builderList[T any] struct {
	list   []*T
	errors []error
	f      *Factory[T]
}

func New[T any](v T) *Factory[T] {
	dataType := reflect.TypeOf(v)

	// TODO: check the edge case
	// 1. tag is comma separated
	// 2. tag is not provided
	// 3. tag is not in the format of "structName,tableName"
	tagToInfo := map[string]tagInfo{}
	for i := 0; i < dataType.NumField(); i++ {
		field := dataType.Field(i)
		tag := field.Tag.Get("factory")
		if tag == "" {
			continue
		}

		parts := strings.Split(tag, ",")
		structName := parts[0]
		tableName := parts[1]

		tagToInfo[structName] = tagInfo{tableName: tableName, fieldName: field.Name}
	}

	last := reflect.New(dataType).Elem().Interface().(T)
	return &Factory[T]{
		dataType:     dataType,
		empty:        last,
		associations: map[string][]interface{}{},
		tagToInfo:    tagToInfo,
		index:        1,
	}
}

// SetConfig sets the configuration for the factory
func (f *Factory[T]) SetConfig(c Config[T]) *Factory[T] {
	f.bluePrint = c.BluePrint
	f.inserter = c.Inserter
	f.db = c.DB
	f.tableName = c.TableName

	return f
}

// SetTrait adds a trait to the factory
func (f *Factory[T]) SetTrait(name string, tr setTraiter[T]) *Factory[T] {
	if f.traits == nil {
		f.traits = map[string]setTraiter[T]{}
	}

	f.traits[name] = tr
	return f
}

// Reset resets the factory to its initial state
func (f *Factory[T]) Reset() {
	f.index = 1
}

// Build builds a value
func (f *Factory[T]) Build() *builder[T] {
	var v T
	if f.bluePrint == nil {
		setNonZeroValues(f.index, &v)
	} else {
		v = f.bluePrint(f.index, v)
	}
	f.index++

	return &builder[T]{
		v:      &v,
		errors: []error{},
		f:      f,
	}
}

// BuildList creates a list of n values
func (f *Factory[T]) BuildList(n int) *builderList[T] {
	list := make([]*T, n)
	errors := []error{}
	if n < 1 {
		errors = append(errors, errBuildListNGreaterThanZero)
		return &builderList[T]{errors: errors}
	}

	for i := 0; i < n; i++ {
		var v T
		if f.bluePrint == nil {
			setNonZeroValues(f.index, &v)
		} else {
			v = f.bluePrint(f.index, v)
		}
		list[i] = &v
		f.index++
	}

	return &builderList[T]{
		list:   list,
		errors: errors,
		f:      f,
	}
}

// Get returns the value
func (b *builder[T]) Get() (T, error) {
	if len(b.errors) > 0 {
		return b.f.empty, genFinalError(b.errors)
	}

	return *b.v, nil
}

// Get returns the list of values
func (b *builderList[T]) Get() ([]T, error) {
	if len(b.errors) > 0 {
		return nil, genFinalError(b.errors)
	}

	output := make([]T, len(b.list))
	for i, v := range b.list {
		output[i] = *v
	}

	return output, nil
}

// Insert inserts the value into the database
func (b *builder[T]) Insert() (T, error) {
	if b.f.db == nil {
		b.errors = append(b.errors, errDBNotProvided)
	}

	if len(b.errors) > 0 {
		return b.f.empty, genFinalError(b.errors)
	}

	// use default inserter
	if b.f.inserter == nil {
		// tableName must provided if not providing Inserter
		if b.f.tableName == "" {
			b.errors = append(b.errors, fmt.Errorf("Insert: %s", errTableNameNotProvided))
			return b.f.empty, genFinalError(b.errors)
		}

		v, err := insertToDB(b.f.db, b.f.tableName, b.v)
		if err != nil {
			b.errors = append(b.errors, err)
			return b.f.empty, genFinalError(b.errors)
		}

		// TODO
		vv, ok := v[0].(*T)
		if !ok {
			b.errors = append(b.errors, fmt.Errorf("Insert: invalid type %T", v[0]))
			return b.f.empty, genFinalError(b.errors)
		}

		return *vv, nil
	}

	// use custom inserter
	v, err := b.f.inserter(b.f.db, *b.v)
	if err != nil {
		b.errors = append(b.errors, err)
		return b.f.empty, genFinalError(b.errors)
	}

	return v, nil
}

// Insert inserts the list of values into the database
func (b *builderList[T]) Insert() ([]T, error) {
	if b.f.db == nil {
		b.errors = append(b.errors, errDBNotProvided)
	}

	if len(b.errors) > 0 {
		return nil, genFinalError(b.errors)
	}

	input := make([]interface{}, len(b.list))
	for i, v := range b.list {
		input[i] = v
	}

	// use default inserter
	if b.f.inserter == nil {
		// tableName must provided if not providing Inserter
		if b.f.tableName == "" {
			b.errors = append(b.errors, fmt.Errorf("Insert: %s", errTableNameNotProvided))
			return nil, genFinalError(b.errors)
		}

		vals, err := insertToDB(b.f.db, b.f.tableName, input...)
		if err != nil {
			b.errors = append(b.errors, err)
			return nil, genFinalError(b.errors)
		}

		// TODO
		output := make([]T, len(vals))
		for i, v := range vals {
			vt, ok := v.(*T)
			if !ok {
				b.errors = append(b.errors, fmt.Errorf("InsertList: invalid type %T", v))
				return nil, genFinalError(b.errors)
			}

			output[i] = *vt
		}

		return output, nil
	}

	// use custom inserter
	output := make([]T, len(b.list))
	for i, v := range b.list {
		val, err := b.f.inserter(b.f.db, *v)
		if err != nil {
			b.errors = append(b.errors, err)
			return nil, genFinalError(b.errors)
		}

		output[i] = val
	}

	return output, nil
}

// Overwrite overwrites the value with the given value
func (b *builder[T]) Overwrite(ow T) *builder[T] {
	if len(b.errors) > 0 {
		return b
	}

	if err := copyValues(b.v, ow); err != nil {
		b.errors = append(b.errors, err)
	}

	return b
}

// Overwrites overwrites the values with the given values
func (b *builderList[T]) Overwrites(ows ...T) *builderList[T] {
	if len(b.errors) > 0 {
		return b
	}

	for i := 0; i < len(ows) && i < len(b.list); i++ {
		if err := copyValues(b.list[i], ows[i]); err != nil {
			b.errors = append(b.errors, err)
			return b
		}
	}

	return b
}

// WithTrait invokes the traiter based on given name
func (b *builder[T]) WithTrait(name string) *builder[T] {
	tr, ok := b.f.traits[name]
	if !ok {
		b.errors = append(b.errors, fmt.Errorf("WithTrait: undefined name %s", name))
		return b
	}

	tr(b.v)

	return b
}

// WithTraits invokes the traiter based on given names
func (b *builderList[T]) WithTraits(names ...string) *builderList[T] {
	for i := 0; i < len(names) && i < len(b.list); i++ {
		tr, ok := b.f.traits[names[i]]
		if !ok {
			b.errors = append(b.errors, fmt.Errorf("WithTraits: undefined name %s", names[i]))
			return b
		}

		tr(b.list[i])
	}

	return b
}

// WihtOne set one association to the factory value
func (b *builder[T]) WithOne(v interface{}) *builder[T] {
	if len(b.errors) > 0 {
		return b
	}

	if err := setAssValue(v, b.f.tagToInfo, b.f.index, "WithOne"); err != nil {
		b.errors = append(b.errors, err)
		return b
	}

	name := reflect.TypeOf(v).Elem().Name()
	b.f.associations[name] = append(b.f.associations[name], v)
	b.f.index++
	return b
}

// WithMany set many associations to the factory value
func (b *builderList[T]) WithMany(values ...interface{}) *builderList[T] {
	if len(b.errors) > 0 {
		return b
	}

	var curValName string
	for _, v := range values {
		if err := setAssValue(v, b.f.tagToInfo, b.f.index, "WithMany"); err != nil {
			b.errors = append(b.errors, err)
			return b
		}

		// check if the provided values are of the same type
		// because we have to make sure all the value is pointer (setAssValue does that for us)
		// before we can use Elem()
		if curValName != "" && curValName != reflect.TypeOf(v).Elem().Name() {
			b.errors = append(b.errors, fmt.Errorf("WithMany: the provided values are not of the same type"))
			return b
		}

		name := reflect.TypeOf(v).Elem().Name()
		b.f.associations[name] = append(b.f.associations[name], v)
		b.f.index++
		curValName = name
	}

	return b
}

// InsertWithAss inserts the value with the associations into the database
func (b *builder[T]) InsertWithAss() (T, []interface{}, error) {
	if len(b.errors) > 0 {
		return b.f.empty, nil, genFinalError(b.errors)
	}

	// generate and insert the associations
	assVals, err := genAndInsertAss(b.f.db, b.f.associations, b.f.tagToInfo)
	if err != nil {
		b.errors = append(b.errors, err)
		return b.f.empty, nil, genFinalError(b.errors)
	}

	// set the connection between the factory value and the associations
	for name, vals := range b.f.associations {
		// use vs[0] because we can make sure InsertWithAss only invoke with Build function
		// which means there's only one factory value
		// so that each associations only allow one value
		fieldName := b.f.tagToInfo[name].fieldName
		if err := setField(b.v, fieldName, vals[0], "InsertWithAss"); err != nil {
			b.errors = append(b.errors, err)
			return b.f.empty, nil, genFinalError(b.errors)
		}
	}

	// insert the factory value
	v, err := b.Insert()
	if err != nil {
		b.errors = append(b.errors, err)
		return b.f.empty, nil, err
	}

	return v, assVals, nil
}

// InsertListWithAss inserts the list of values with the associations into the database
func (b *builderList[T]) InsertWithAss() ([]T, []interface{}, error) {
	if len(b.errors) > 0 {
		return nil, nil, genFinalError(b.errors)
	}

	// generate and insert
	assVals, err := genAndInsertAss(b.f.db, b.f.associations, b.f.tagToInfo)
	if err != nil {
		b.errors = append(b.errors, err)
		return nil, nil, genFinalError(b.errors)
	}

	// set the connection between the factory value and the associations
	cachePrev := map[string]interface{}{}
	for i, l := range b.list {
		for name, vs := range b.f.associations {
			var v interface{}
			if i >= len(vs) {
				v = cachePrev[name]
			} else {
				v = vs[i]
				cachePrev[name] = vs[i]
			}

			fieldName := b.f.tagToInfo[name].fieldName
			if err := setField(l, fieldName, v, "InsertWithAss"); err != nil {
				b.errors = append(b.errors, err)
				return nil, nil, genFinalError(b.errors)
			}
		}
	}

	// insert the factory value
	v, err := b.Insert()
	if err != nil {
		b.errors = append(b.errors, err)
		return nil, nil, genFinalError(b.errors)
	}

	return v, assVals, nil
}

// setField sets the value to the name field of the target
func setField(target interface{}, name string, source interface{}, sourceFn string) error {
	targetField := reflect.ValueOf(target).Elem().FieldByName(name)
	if !targetField.IsValid() {
		return fmt.Errorf("%s: field %s is not found", sourceFn, name)
	}

	if !targetField.CanSet() {
		return fmt.Errorf("%s: field %s can not be set", sourceFn, name)
	}

	sourceIDField := reflect.ValueOf(source).Elem().FieldByName("ID")
	if !sourceIDField.IsValid() {
		return fmt.Errorf("%s: source field ID is not found", sourceFn)
	}

	sourceIDKind := sourceIDField.Kind()
	if sourceIDKind != reflect.Int &&
		sourceIDKind != reflect.Int64 &&
		sourceIDKind != reflect.Int32 &&
		sourceIDKind != reflect.Int16 &&
		sourceIDKind != reflect.Int8 &&
		sourceIDKind != reflect.Uint &&
		sourceIDKind != reflect.Uint64 &&
		sourceIDKind != reflect.Uint32 &&
		sourceIDKind != reflect.Uint16 &&
		sourceIDKind != reflect.Uint8 {
		return fmt.Errorf("%s: source field ID is not an integer", sourceFn)
	}

	// TODO: What if targetField is int, but sourceIDField is uint?
	switch sourceIDField.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		targetField.SetInt(sourceIDField.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		targetField.SetUint(sourceIDField.Uint())
	}

	return nil
}

// genAndInsertAss inserts the associations value into the database and returns with the inserted values
func genAndInsertAss(db *sql.DB, associations map[string][]interface{}, tagToInfo map[string]tagInfo) ([]interface{}, error) {
	if len(associations) == 0 {
		return nil, errInsertAssWithoutAss
	}

	result := []interface{}{}
	for name, vals := range associations {
		tableName := tagToInfo[name].tableName

		v, err := insertToDB(db, tableName, vals...)
		if err != nil {
			return nil, err
		}

		result = append(result, v...)
	}

	return result, nil
}

// setAssValue sets the value to the associations value
func setAssValue(v interface{}, tagToInfo map[string]tagInfo, index int, sourceFn string) error {
	typeOfV := reflect.TypeOf(v)

	// check if it's a pointer
	if typeOfV.Kind() != reflect.Ptr {
		name := typeOfV.Name()
		return fmt.Errorf("%s: type %s, value %v is not a pointer", sourceFn, name, v)
	}

	name := typeOfV.Elem().Name()
	// check if it's a pointer to a struct
	if typeOfV.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("%s: type %s, value %v is not a pointer to a struct", sourceFn, name, v)
	}

	// check if it's existed in tagToInfo
	if _, ok := tagToInfo[name]; !ok {
		return fmt.Errorf("%s: type %s, value %v is not found at tag", sourceFn, name, v)
	}

	setNonZeroValues(index, v)
	return nil
}

// Insert inserts the values into the database, values must be a pointer to a struct
func insertToDB(db *sql.DB, tableName string, values ...interface{}) ([]interface{}, error) {
	rawStmt, fieldValues := prepareStmtAndVals(tableName, values)

	stmt, err := db.Prepare(rawStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result := make([]interface{}, len(fieldValues))
	for i, vals := range fieldValues {
		fmt.Println("insertToDB: ", rawStmt, vals)
		res, err := tx.Stmt(stmt).Exec(vals...)
		if err != nil {
			return nil, err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}

		val := values[i]
		v := reflect.ValueOf(val).Elem()
		idField := v.FieldByName("ID")
		switch idField.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			idField.SetInt(id)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			idField.SetUint(uint64(id))
		}

		result[i] = val
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

// prepareStmtAndVals prepares the SQL insert statement and the values to be inserted
// it's the helper function for insertToDB
func prepareStmtAndVals(tableName string, values []interface{}) (string, [][]interface{}) {
	fieldNames := []string{}
	placeholders := []string{}
	fieldValues := [][]interface{}{}

	for index, val := range values {
		val := reflect.ValueOf(val).Elem()
		vals := []interface{}{}

		for i := 0; i < val.NumField(); i++ {
			n := val.Type().Field(i).Name
			if n == "ID" {
				continue
			}

			vals = append(vals, val.Field(i).Interface())

			if index == 0 {
				fieldNames = append(fieldNames, camelToSnake(n))
				placeholders = append(placeholders, "?")
			}
		}

		fieldValues = append(fieldValues, vals)
	}

	// Construct the SQL insert statement
	fns := strings.Join(fieldNames, ", ")
	phs := strings.Join(placeholders, ", ")
	rawStmt := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, fns, phs)

	return rawStmt, fieldValues
}
