package database

import (
	"database/sql"
	"log"
	"umbandario-go/models"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(filepath string) {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatalf("Erro ao abrir o banco de dados: %v", err)
	}

	createAudioTableSQL := `
	CREATE TABLE IF NOT EXISTS audio_files (
		"id" TEXT NOT NULL PRIMARY KEY,
		"filename" TEXT UNIQUE,
		"filetype" TEXT,
		"path" TEXT,
		"line_id" TEXT,
		"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(line_id) REFERENCES lines(id) ON DELETE SET NULL -- Regra da chave estrangeira
	);`
	createAudioTriggerSQL := `
	CREATE TRIGGER IF NOT EXISTS update_audio_files_updated_at
	AFTER UPDATE ON audio_files
	FOR EACH ROW
	BEGIN
		UPDATE audio_files SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
	END;`

	createLinesTableSQL := `
	CREATE TABLE IF NOT EXISTS lines (
		"id" TEXT NOT NULL PRIMARY KEY,
		"name" TEXT NOT NULL UNIQUE,
		"created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		"updated_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	createLinesTriggerSQL := `
	CREATE TRIGGER IF NOT EXISTS update_lines_updated_at
	AFTER UPDATE ON lines
	FOR EACH ROW
	BEGIN
		UPDATE lines SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
	END;`

	execSQL(createLinesTableSQL, "tabela 'lines'")
	execSQL(createLinesTriggerSQL, "trigger 'update_lines_updated_at'")
	execSQL(createAudioTableSQL, "tabela 'audio_files'")
	execSQL(createAudioTriggerSQL, "trigger 'update_audio_files_updated_at'")

	log.Println("Banco de dados inicializado com as tabelas 'audio_files' e 'lines'.")
}

func execSQL(sql, name string) {
	if _, err := DB.Exec(sql); err != nil {
		log.Fatalf("Erro ao criar %s: %v", name, err)
	}
}

func CreateAudioFile(audio models.AudioFile) (models.AudioFile, error) {
	_, err := DB.Exec("INSERT INTO audio_files (id, filename, filetype, path, line_id) VALUES (?, ?, ?, ?, ?)",
		audio.ID, audio.Filename, audio.Filetype, audio.Path, audio.LineID)
	if err != nil {
		return models.AudioFile{}, err
	}
	createdAudio, err := GetAudioFileByID(audio.ID)
	if err != nil {
		return models.AudioFile{}, err
	}
	return *createdAudio, nil
}

func GetAudioFileByID(id string) (*models.AudioFile, error) {
	row := DB.QueryRow("SELECT id, filename, filetype, path, line_id, created_at, updated_at FROM audio_files WHERE id = ?", id)
	var audio models.AudioFile
	if err := row.Scan(&audio.ID, &audio.Filename, &audio.Filetype, &audio.Path, &audio.LineID, &audio.CreatedAt, &audio.UpdatedAt); err != nil {
		return nil, err
	}
	return &audio, nil
}

func GetAllAudioFiles() (*models.AudioFileList, error) {
	rows, err := DB.Query("SELECT id, filename, filetype, path, line_id, created_at, updated_at FROM audio_files ORDER BY filename")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var audioFiles = make([]models.AudioFile, 0)
	for rows.Next() {
		var audio models.AudioFile
		if err := rows.Scan(&audio.ID, &audio.Filename, &audio.Filetype, &audio.Path, &audio.LineID, &audio.CreatedAt, &audio.UpdatedAt); err != nil {
			return nil, err
		}
		audioFiles = append(audioFiles, audio)
	}

	var getAllAudio models.AudioFileList
	getAllAudio.Dados = audioFiles
	return &getAllAudio, nil
}

func DeleteAudioFile(id string) (int64, error) {
	result, err := DB.Exec("DELETE FROM audio_files WHERE id = ?", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func CreateLine(name string) (*models.Line, error) {
	id := uuid.NewString()
	_, err := DB.Exec("INSERT INTO lines (id, name) VALUES (?, ?)", id, name)
	if err != nil {
		return nil, err
	}
	return GetLineByID(id)
}

func GetAllLines() (*models.LineList, error) {
	rows, err := DB.Query("SELECT id, name, created_at, updated_at FROM lines ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lines = make([]models.Line, 0)
	for rows.Next() {
		var line models.Line
		if err := rows.Scan(&line.ID, &line.Name, &line.CreatedAt, &line.UpdatedAt); err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	var getAllLines models.LineList
	getAllLines.Dados = lines
	return &getAllLines, nil
}

func GetLineByName(name string) (*models.Line, error) {
	row := DB.QueryRow("SELECT id, name, created_at, updated_at FROM lines WHERE name = ?", name)
	var line models.Line
	if err := row.Scan(&line.ID, &line.Name, &line.CreatedAt, &line.UpdatedAt); err != nil {
		return nil, err
	}
	return &line, nil
}

func GetLineByID(id string) (*models.Line, error) {
	row := DB.QueryRow("SELECT id, name, created_at, updated_at FROM lines WHERE id = ?", id)
	var line models.Line
	if err := row.Scan(&line.ID, &line.Name, &line.CreatedAt, &line.UpdatedAt); err != nil {
		return nil, err
	}

	return &line, nil
}

func UtilScanAndStoreAudios(dir string) {}
