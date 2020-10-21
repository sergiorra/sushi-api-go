package inmem

import (
	"fmt"
	"sync"

	sushi "github.com/sergiorra/sushi-api-go/pkg"
)

type sushiRepository struct {
	mtx     sync.RWMutex
	sushis 	map[string]sushi.Sushi
}

func NewRepository(sushis map[string]sushi.Sushi) sushi.Repository {
	if sushis == nil {
		sushis = make(map[string]sushi.Sushi)
	}

	return &sushiRepository{
		sushis: sushis,
	}
}

func (r *sushiRepository) CreateSushi(s *sushi.Sushi) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if err := r.checkIfExists(s.ID); err != nil {
		return err
	}
	r.sushis[s.ID] = *s
	return nil
}

func (r *sushiRepository) GetSushis() ([]sushi.Sushi, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	values := make([]sushi.Sushi, 0, len(r.sushis))
	for _, value := range r.sushis {
		values = append(values, value)
	}
	return values, nil
}

func (r *sushiRepository) DeleteSushi(ID string) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	delete(r.sushis, ID)

	return nil
}

func (r *sushiRepository) UpdateSushi(ID string, s *sushi.Sushi) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.sushis[ID] = *s
	return nil
}

func (r *sushiRepository) GetSushiByID(ID string) (*sushi.Sushi, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	for _, v := range r.sushis {
		if v.ID == ID {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("The ID %s does not exist", ID)
}

func (r *sushiRepository) checkIfExists(ID string) error {
	for _, v := range r.sushis {
		if v.ID == ID {
			return fmt.Errorf("The sushi %s already exists", ID)
		}
	}

	return nil
}