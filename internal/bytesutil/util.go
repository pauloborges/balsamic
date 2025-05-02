package bytesutil

import (
	"bytes"
)

// Buffer is a wrapper around bytes.Buffer that provides additional
// functionality.
type Buffer struct {
	bytes.Buffer
}

// WriteWithSuffix writes the bytes from p to the buffer, followed by the
// bytes from suffix. It only writes the suffix if p is not empty.
func (b *Buffer) WriteWithSuffix(p []byte, suffix string) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	n, err := b.Buffer.Write(p)
	if err != nil {
		return n, err
	}

	ns, err := b.Buffer.WriteString(suffix)
	return n + ns, err
}

// WriteWithPrefix writes the bytes from p to the buffer, preceded by the
// bytes from prefix. It only writes the prefix if p is not empty.
func (b *Buffer) WriteWithPrefix(prefix string, p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	np, err := b.Buffer.WriteString(prefix)
	if err != nil {
		return np, err
	}

	n, err := b.Buffer.Write(p)
	return np + n, err
}

func (b *Buffer) WriteWithPrefixSuffix(prefix string, p []byte, suffix string) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	np, err := b.Buffer.WriteString(prefix)
	if err != nil {
		return np, err
	}

	n, err := b.Buffer.Write(p)
	if err != nil {
		return np + n, err
	}

	ns, err := b.Buffer.WriteString(suffix)
	return np + n + ns, err
}
