package efactory

// Database is responsible for inserting data into the database
type Database interface {
	// insert inserts a single data into the database
	insert(inserParams) (interface{}, error)

	// insertList inserts a list of data into the database
	insertList(inserListParams) ([]interface{}, error)
}
