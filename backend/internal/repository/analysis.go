package repository

import (
	"database/sql"
	"fmt"

	"twin-os/backend/internal/model"
)

// AnalysisRepository 分析结果仓储
type AnalysisRepository struct {
	db *sql.DB
}

// NewAnalysisRepository 创建分析结果仓储
func NewAnalysisRepository(db *sql.DB) *AnalysisRepository {
	return &AnalysisRepository{db: db}
}

// Create 创建分析结果
func (r *AnalysisRepository) Create(analysis *model.AnalysisResult) error {
	query := `
		INSERT INTO analysis_results (
			type, message_ids, content, metadata, confidence, created_at, processed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query,
		analysis.Type,
		analysis.MessageIDs,
		analysis.Content,
		analysis.Metadata,
		analysis.Confidence,
		analysis.CreatedAt,
		analysis.ProcessedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create analysis: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	analysis.ID = id
	return nil
}

// GetByType 根据类型获取最新的分析结果
func (r *AnalysisRepository) GetByType(analysisType string) (*model.AnalysisResult, error) {
	query := `
		SELECT id, type, message_ids, content, metadata, confidence, created_at, processed_at
		FROM analysis_results
		WHERE type = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	row := r.db.QueryRow(query, analysisType)
	analysis := &model.AnalysisResult{}

	err := row.Scan(
		&analysis.ID,
		&analysis.Type,
		&analysis.MessageIDs,
		&analysis.Content,
		&analysis.Metadata,
		&analysis.Confidence,
		&analysis.CreatedAt,
		&analysis.ProcessedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("analysis not found")
		}
		return nil, fmt.Errorf("failed to get analysis: %w", err)
	}

	return analysis, nil
}