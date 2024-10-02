package store

import (
	// 필요한 패키지들을 임포트합니다.
	"Bangseungjae/go-todo-app/clock"
	"Bangseungjae/go-todo-app/config"
	"context"
	"database/sql"
	"fmt"
	"time"

	// MySQL 드라이버를 익명 임포트하여 초기화 사이드 이펙트를 발생시킵니다.
	_ "github.com/go-sql-driver/mysql"
	// 확장된 SQL 기능을 제공하는 sqlx 패키지를 임포트합니다.
	"github.com/jmoiron/sqlx"
)

// New 함수는 데이터베이스 연결을 생성하고 반환합니다.
// 또한, 데이터베이스 연결을 닫는 클린업 함수와 에러를 반환합니다.
func New(ctx context.Context, cfg *config.Config) (*sqlx.DB, func(), error) {
	// 데이터베이스 연결을 초기화합니다.
	db, err := sql.Open("mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.DBUser, cfg.DBPassword, // 사용자 이름과 비밀번호
			cfg.DBHost, cfg.DBPort, // 호스트와 포트
			cfg.DBName, // 데이터베이스 이름
		),
	)
	if err != nil {
		// 데이터베이스 초기화에 실패하면 에러를 반환합니다.
		return nil, func() {}, err
	}
	// Open 함수는 실제로 데이터베이스에 연결하지 않으므로 Ping을 통해 연결을 확인합니다.
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		// Ping에 실패하면 데이터베이스를 닫고 에러를 반환합니다.
		return nil, func() { _ = db.Close() }, err
	}
	// sqlx 패키지의 NewDb 함수를 사용하여 *sql.DB를 *sqlx.DB로 래핑합니다.
	xdb := sqlx.NewDb(db, "mysql")
	// 데이터베이스 연결과 클린업 함수를 반환합니다.
	return xdb, func() { _ = db.Close() }, nil
}

// Beginner 인터페이스는 트랜잭션을 시작하는 메서드를 정의합니다.
type Beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// Preparer 인터페이스는 쿼리를 사전에 준비하는 메서드를 정의합니다.
type Preparer interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}

// Execer 인터페이스는 SQL 쿼리를 실행하는 메서드를 정의합니다.
type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

// Queryer 인터페이스는 데이터베이스에서 쿼리를 수행하고 결과를 가져오는 메서드를 정의합니다.
type Queryer interface {
	Preparer // Preparer 인터페이스를 포함하여 PreparexContext 메서드를 상속합니다.
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...any) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error
}

var (
	// 인터페이스들이 예상대로 구현되었는지 컴파일 타임에 확인합니다.
	_ Beginner = (*sqlx.DB)(nil) // sqlx.DB가 Beginner 인터페이스를 구현하는지 확인
	_ Preparer = (*sqlx.DB)(nil) // sqlx.DB가 Preparer 인터페이스를 구현하는지 확인
	_ Queryer  = (*sqlx.DB)(nil) // sqlx.DB가 Queryer 인터페이스를 구현하는지 확인
	_ Execer   = (*sqlx.DB)(nil) // sqlx.DB가 Execer 인터페이스를 구현하는지 확인
	_ Execer   = (*sqlx.Tx)(nil) // sqlx.Tx가 Execer 인터페이스를 구현하는지 확인
)

// Repository 구조체는 데이터베이스 작업에 필요한 의존성을 포함합니다.
type Repository struct {
	Clocker clock.Clocker // 현재 시간을 얻기 위한 Clocker 인터페이스
}
