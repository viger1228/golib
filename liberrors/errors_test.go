package liberrors

import (
	"context"
	"errors"
	"log"
	"testing"
)

func TestRun(t *testing.T) {
	var err error
	ctx := context.Background()

	err = Func1(ctx)
	if err != nil {
		log.Printf(ErrorFormat(err), err.Error())
	}
	err = Func4(ctx)
	if err != nil {
		log.Printf(ErrorFormat(err), err.Error())
	}
}

func BenchmarkRun(b *testing.B) {
	//與性能測試無關的code
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//測試的code //var err error
	}
	b.StopTimer()
	//與性能測試無關的code
}

func Func1(ctx context.Context) (err error) {
	err = Func2(ctx)
	if err != nil {
		log.Printf(ErrorFormat(err), err.Error())
		return err
	}
	return nil
}

func Func2(ctx context.Context) (err error) {
	err = Func3(ctx)
	if err != nil {
		log.Printf(ErrorFormat(err), err.Error())
		return err
	}
	return nil
}

func Func3(ctx context.Context) (err error) {
	err = MyError(ctx)
	if err != nil {
		log.Printf(ErrorFormat(err), err.Error())
		return err
	}
	return nil
}

func Func4(ctx context.Context) (err error) {
	err = ThirdPartyError(ctx)
	if err != nil {
		log.Printf(ErrorFormat(err), err.Error())
		return err
	}
	return nil
}

func MyError(ctx context.Context) (err error) {
	err = New(ctx, "my error")
	AppendLink(err, "Hello")
	log.Printf(ErrorFormat(err), err.Error())
	return err
}

func ThirdPartyError(ctx context.Context) (err error) {
	err = errors.New("third party error")
	err = Load(ctx, err)
	log.Printf(ErrorFormat(err), err.Error())
	return err
}
