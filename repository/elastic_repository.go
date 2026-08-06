package repository

type Log struct {
	ID      string `json:"-"`
	Message string `json:"message"`
	Level   string `json:"level"`
}

/*
type LogElasticRepository struct {
	es *elasticsearch.Client
}

func NewLogElasticRepository(es *elasticsearch.Client) *LogElasticRepository {
	return &LogElasticRepository{es: es}
}

func (r *LogElasticRepository) SearchLogs(ctx context.Context, keyword string, page, limit int) ([]Log, int64, error) {
	from := (page - 1) * limit

	// query DSL dalam bentuk map, nanti di-encode ke JSON
	query := map[string]interface{}{
		"from": from,
		"size": limit,
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"message": keyword,
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, err
	}

	res, err := r.es.Search(
		r.es.Search.WithContext(ctx),
		r.es.Search.WithIndex("logs"), // nama index
		r.es.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	var result struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				ID     string          `json:"_id"`
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	logs := make([]Log, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		var l Log
		if err := json.Unmarshal(hit.Source, &l); err != nil {
			continue
		}
		l.ID = hit.ID
		logs = append(logs, l)
	}

	return logs, result.Hits.Total.Value, nil
}
*/
