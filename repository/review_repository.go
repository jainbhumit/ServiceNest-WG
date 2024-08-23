package repository

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"serviceNest/model"
	"sync"
)

type ReviewRepository struct {
	filePath string
	mutex    sync.Mutex
}

func NewReviewRepository(filePath string) *ReviewRepository {
	return &ReviewRepository{
		filePath: filePath,
	}
}

// Load all reviews from the JSON file
func (r *ReviewRepository) loadReviews() ([]model.Review, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	file, err := os.Open(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Review{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var reviews []model.Review
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&reviews)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

// Save all reviews to the JSON file
func (r *ReviewRepository) saveReviews(reviews []model.Review) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	data, err := json.MarshalIndent(reviews, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(r.filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Get a review by ID
func (r *ReviewRepository) GetReviewByID(reviewID string) (*model.Review, error) {
	reviews, err := r.loadReviews()
	if err != nil {
		return nil, err
	}

	for _, review := range reviews {
		if review.ID == reviewID {
			return &review, nil
		}
	}
	return nil, errors.New("review not found")
}

// Save a new review
func (r *ReviewRepository) SaveReview(review model.Review) error {
	reviews, err := r.loadReviews()
	if err != nil {
		return err
	}

	reviews = append(reviews, review)
	return r.saveReviews(reviews)
}

// Update an existing review
func (r *ReviewRepository) UpdateReview(updatedReview model.Review) error {
	reviews, err := r.loadReviews()
	if err != nil {
		return err
	}

	for i, review := range reviews {
		if review.ID == updatedReview.ID {
			reviews[i] = updatedReview
			return r.saveReviews(reviews)
		}
	}
	return errors.New("review not found")
}

// Delete a review by ID
func (r *ReviewRepository) DeleteReview(reviewID string) error {
	reviews, err := r.loadReviews()
	if err != nil {
		return err
	}

	for i, review := range reviews {
		if review.ID == reviewID {
			reviews = append(reviews[:i], reviews[i+1:]...)
			return r.saveReviews(reviews)
		}
	}
	return errors.New("review not found")
}

// Get all reviews
func (r *ReviewRepository) GetAllReviews() ([]model.Review, error) {
	return r.loadReviews()
}
