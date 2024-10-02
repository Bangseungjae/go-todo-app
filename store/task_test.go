package store

import (
	"Bangseungjae/go-todo-app/clock"
	"Bangseungjae/go-todo-app/entity"
	"Bangseungjae/go-todo-app/testutil"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	"testing"
)

func newMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	xdb := sqlx.NewDb(db, "mysql")
	return xdb, mock
}

func TestRepository_AddTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}
	var wantID int64 = 20
	okTask := &entity.Task{
		Title:    "ok task",
		Status:   "todo",
		Created:  c.Now(),
		Modified: c.Now(),
	}

	xdb, mock := newMockDB(t)
	sqlQuery := "INSERT INTO task (title, status, created, modified) VALUES (?, ?, ?, ?)"
	mock.ExpectExec(sqlQuery).
		WithArgs(okTask.Title, okTask.Status, okTask.Created, okTask.Modified).
		WillReturnResult(sqlmock.NewResult(wantID, 1))

	r := &Repository{Clocker: c}
	if err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Fatalf("want no error, but got %v", err)
	}
}

func TestRepository_ListTasks(t *testing.T) {
	ctx := context.Background()
	// entity.Task를 작성하는 다른 테스트 케이스와 섞이면 테스트가 실패한다.
	// 이를 위해 트랜잭션을 적용해서 테스트 케이스 내로 한정된 테이블 상태를 만든다.
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	// 이 테스트 케이스가 끝나면 원래 상태로 되돌린다.
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}
	wants := prepareTasks(ctx, t, tx)

	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}

func prepareTasks(ctx context.Context, t *testing.T, con Execer) entity.Tasks {
	t.Helper()
	// 깨끗한 상태로 정리해둔다.
	if _, err := con.ExecContext(ctx, "DELETE FROM task;"); err != nil {
		t.Logf("failed to initialize task: %v", err)
	}
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			Title: "want task 1", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		{
			Title: "want task 2", Status: "todo",
			Created: c.Now(), Modified: c.Now(),
		},
		{
			Title: "want task 3", Status: "done",
			Created: c.Now(), Modified: c.Now(),
		},
	}
	result, err := con.ExecContext(ctx,
		`INSERT INTO task (title, status, created, modified)
			VALUES
			    (?, ?, ?, ?),
			    (?, ?, ?, ?),
			    (?, ?, ?, ?);`,
		wants[0].Title, wants[0].Status, wants[0].Created, wants[0].Modified,
		wants[1].Title, wants[1].Status, wants[1].Created, wants[1].Modified,
		wants[2].Title, wants[2].Status, wants[2].Created, wants[2].Modified,
	)
	if err != nil {
		t.Fatal(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	wants[0].ID = entity.TaskID(id)
	wants[1].ID = entity.TaskID(id + 1)
	wants[2].ID = entity.TaskID(id + 2)
	return wants
}
