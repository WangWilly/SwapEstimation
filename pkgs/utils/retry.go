package utils

import "time"

////////////////////////////////////////////////////////////////////////////////

func retry(attempts int, sleep time.Duration, fn func(i int) error) error {
	if err := fn(attempts); err != nil {
		if attempts--; attempts > 0 {
			time.Sleep(sleep)
			return retry(attempts, sleep, fn)
		}
		return err
	}
	return nil
}
