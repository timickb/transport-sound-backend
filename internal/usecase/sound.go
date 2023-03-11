package usecase

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/timickb/transport-sound/internal/domain"
	"math/rand"
	"time"
)

type SoundUseCase struct {
	r Repository
}

func NewSoundUseCase(r Repository) *SoundUseCase {
	return &SoundUseCase{r: r}
}

func (u *SoundUseCase) GetSoundById(id string) (*domain.Sound, error) {
	sound, err := u.r.GetSoundById(id)
	if err != nil {
		return nil, fmt.Errorf("get sound by id err: %w", err)
	}

	tags, err := u.r.GetTagsForSound(sound.Id)
	if err != nil {
		return nil, fmt.Errorf("get sound by id err: %w", err)
	}

	sound.Tags = tags
	return sound, nil
}

func (u *SoundUseCase) GetAllSounds() ([]*domain.Sound, error) {
	sounds, err := u.r.GetAllSounds()
	if err != nil {
		return nil, err
	}

	return sounds, nil
}

func (u *SoundUseCase) CreateSound(ctx UserContext, s *domain.Sound, tid []string) (string, error) {
	// Check each tag for existing.
	for _, tagId := range tid {
		_, err := u.r.GetTagById(tagId)
		if err != nil {
			return "", fmt.Errorf("err create sound: %w", err)
		}
	}

	// Create the sound.
	if s.VehicleId == "" {
		s.VehicleId = "default"
	}
	s.Id = uuid.NewString()
	err := u.r.CreateSound(s)
	if err != nil {
		return "", fmt.Errorf("err create sound: %w", err)
	}

	// Bind tags to the created sound.
	for _, tagId := range tid {
		if err := u.r.AddTagToSound(s.Id, tagId); err != nil {
			return "", fmt.Errorf("err create sound: %w", err)
		}
	}

	return s.Id, nil
}

func (u *SoundUseCase) GetRandomSounds(limit int) ([]*domain.Sound, error) {
	sounds, err := u.r.GetAllSounds()
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(sounds), func(i, j int) {
		sounds[i], sounds[j] = sounds[j], sounds[i]
	})

	if len(sounds) <= limit {
		return sounds, nil
	}

	return sounds[:limit], nil
}
