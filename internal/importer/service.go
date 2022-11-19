package importer

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/mailru/easyjson/jwriter"

	"github.com/burnb/duocard/internal/configs"
)

type Service struct {
	cfg    *configs.App
	client *http.Client
}

func New(cfg *configs.App) *Service {
	return &Service{
		cfg:    cfg,
		client: http.DefaultClient,
	}
}

func (s *Service) Run() error {
	file, err := os.Open(s.cfg.FilePath)
	if err != nil {
		return fmt.Errorf("unable to open input file %s: %v", s.cfg.FilePath, err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Panic(err)
		}
	}()

	var cards []*Card
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		parts := strings.Split(row, s.cfg.Separator)

		if len(parts) != 3 {
			return errors.New("wrong input file format or separator")
		}

		cards = append(cards, &Card{Front: parts[1], Back: parts[0], Hint: parts[2]})
	}
	s.send(cards)

	return nil
}

func (s *Service) send(cards []*Card) {
	var wg sync.WaitGroup
	wg.Add(len(cards))

	for _, card := range cards {
		go func(card *Card, wg *sync.WaitGroup) {
			defer wg.Done()

			req, err := s.prepareRequest(card)
			if err != nil {
				log.Printf("Unable to prepare request\nCard: %v\nError: %v", card, err)
			}

			res, err := s.client.Do(req)
			if err != nil {
				log.Printf("Unable to do request\nCard: %v\nCode: %d\nError: %v", card, res.StatusCode, err)
			}
		}(card, &wg)
	}

	wg.Wait()
}

func (s *Service) prepareRequest(card *Card) (*http.Request, error) {
	w := &jwriter.Writer{}
	body := &Request{
		Query: "mutation cardCreateMutation(\n  $deckId: ID!\n  $front: String!\n  $back: String!\n  $langBack: String\n  $hint: String\n  $sCardId: ID\n  $sBackId: ID\n  $sourceId: ID\n  $fromSharing: Boolean\n  $deleted: Boolean\n  $svg: FlatIconEdit\n) {\n  cardCreate(deckId: $deckId, front: $front, back: $back, langBack: $langBack, hint: $hint, sCardId: $sCardId, sBackId: $sBackId, sourceId: $sourceId, fromSharing: $fromSharing, deleted: $deleted, svg: $svg) {\n    viewer {\n      xps {\n        cardsToday\n        days {\n          xps\n          added\n          reviewed\n          day\n          streaked\n        }\n      }\n      id\n    }\n    deck {\n      ...loopQuery_stats\n      id\n    }\n    sCard {\n      isInMyDeck(deckId: $deckId)\n      id\n    }\n    duplicatedCard {\n      id\n      front\n      back\n    }\n  }\n}\n\nfragment loopQuery_stats on Deck {\n  stats {\n    known\n    unknown\n    completed\n    total\n    under10\n  }\n}\n",
		Variables: &CreateCardVariables{
			DeckId:   s.cfg.DeckId,
			Front:    card.Front,
			Back:     card.Back,
			LangBack: s.cfg.NativeLanguage,
			Hint:     card.Hint,
		},
	}
	body.MarshalEasyJSON(w)

	req, err := http.NewRequest("POST", "https://api.duocards.com/graphql?cardCreateMutation", w.Buffer.ReadCloser())
	if err != nil {
		return nil, fmt.Errorf("unable to get new http request %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+s.cfg.ApiToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", defaultUserAgent)

	return req, nil
}
