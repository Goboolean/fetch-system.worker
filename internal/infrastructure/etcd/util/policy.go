package etcdutil

// ETCD UTIL POLICY DOCUMENT



// Policy for etcd object:

// A. etcd key format
// First compartment means the type of data.
// It starts with / which means it folows described policy.
// Second compartment means unique id
// Third and more compartment means the field of data.
// If any fourth and more compartment exists, it is a record of third compartment which is an object

// B. etcd model format
// Every struct must have a function named Name() string which returns the name of the struct.
// It is used to identify the type of data.
// Every struct record that needs to be recognized by this package should have a tag named "etcd".


type Model interface {
	Name() string
}


// Policy for etcd semaphore

// A. etcd key format
// Key-value pair for specific usecase does not starts with /, should starts with its function directly.
// Semaphore starts with @, its name follows.


func Semaphore(key string) string {
	return "@" + key
}