package metadata

import (
	"sync"
	"github.com/google/uuid"
)



type Manager struct {
	m sync.Mutex

	platform string
	workerId string
	
	productList map[string]struct{}
	newProductList map[string]struct{}
}

func New(platform string) *Manager {
	uuid := uuid.New()

	return &Manager{
		platform: platform,
		workerId: uuid.String(),
		productList: make(map[string]struct{}),
		newProductList: make(map[string]struct{}),
	}
}


func (m *Manager) GetPlatform() string {
	return m.platform
}

func (m *Manager) GetWorkerId() string {
	return m.workerId
}


func (m *Manager) AddProduct(product string) {
	m.m.Lock()
	defer m.m.Unlock()

	if m.productList == nil {
		m.productList = make(map[string]struct{})
	}

	m.productList[product] = struct{}{}
}

func (m *Manager) ProductExist(product string) bool {
	m.m.Lock()
	defer m.m.Unlock()

	_, exist := m.productList[product]
	return exist
}

func (m *Manager) GetProductList() []string {
	m.m.Lock()
	defer m.m.Unlock()

	productList := make([]string, 0, len(m.productList))
	for product := range m.productList {
		productList = append(productList, product)
	}

	return productList
}


func (m *Manager) AddNewProduct(product string) {
	m.m.Lock()
	defer m.m.Unlock()

	if m.newProductList == nil {
		m.newProductList = make(map[string]struct{})
	}

	m.newProductList[product] = struct{}{}
}

func (m *Manager) RemoveNewProduct(product string) bool {
	m.m.Lock()
	defer m.m.Unlock()

	if m.newProductList == nil {
		return false
	}

	_, exist := m.newProductList[product]
	if !exist {
		return false
	}

	delete(m.newProductList, product)
	return true
}

func (m *Manager) GetNewProductList() []string {
	m.m.Lock()
	defer m.m.Unlock()

	productList := make([]string, 0, len(m.newProductList))
	for product := range m.newProductList {
		productList = append(productList, product)
	}

	return productList
}