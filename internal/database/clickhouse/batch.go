package clickhouse

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"
)

var (
	ErrSliceIsNil      = errors.New("data slice is nil")
	ErrIsNotSlice      = errors.New("data is not a slice")
	ErrItemIsNotStruct = errors.New("batch slice item is not a struct")
)

// nolint:gocyclo
func (c *Connection) SendBatch(ctx context.Context, query string, data interface{}) error {
	if data == nil {
		return ErrSliceIsNil
	}

	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		return ErrIsNotSlice
	}

	itemsCount := value.Len()
	if itemsCount == 0 {
		return nil
	}

	item := value.Index(0)
	if item.Kind() != reflect.Struct {
		return ErrItemIsNotStruct
	}

	columnCount := 0
	var columnIndexes []int
	var columns []interface{}
	var columnTypes []string

	for i := 0; i < item.NumField(); i++ {
		field := item.Type().Field(i)
		tagValue, ok := field.Tag.Lookup("ch")
		if ok && tagValue == "-" {
			continue
		}
		columnCount++
		columnIndexes = append(columnIndexes, i)

		switch field.Type.Name() {
		case "Time":
			columnTypes = append(columnTypes, "Time")
			columns = append(columns, make([]time.Time, itemsCount))
		case "string":
			columnTypes = append(columnTypes, "string")
			columns = append(columns, make([]string, itemsCount))
		case "int8":
			columnTypes = append(columnTypes, "int8")
			columns = append(columns, make([]int8, itemsCount))
		case "int16":
			columnTypes = append(columnTypes, "int16")
			columns = append(columns, make([]int16, itemsCount))
		case "int32":
			columnTypes = append(columnTypes, "int32")
			columns = append(columns, make([]int32, itemsCount))
		case "int64", "int":
			columnTypes = append(columnTypes, "int64")
			columns = append(columns, make([]int64, itemsCount))
		case "uint8":
			columnTypes = append(columnTypes, "uint8")
			columns = append(columns, make([]uint8, itemsCount))
		case "uint16":
			columnTypes = append(columnTypes, "uint16")
			columns = append(columns, make([]uint16, itemsCount))
		case "uint32":
			columnTypes = append(columnTypes, "uint32")
			columns = append(columns, make([]uint32, itemsCount))
		case "uint64":
			columnTypes = append(columnTypes, "uint64")
			columns = append(columns, make([]uint64, itemsCount))
		case "float32":
			columnTypes = append(columnTypes, "float32")
			columns = append(columns, make([]float32, itemsCount))
		case "float64":
			columnTypes = append(columnTypes, "float64")
			columns = append(columns, make([]float64, itemsCount))
		default:
			return fmt.Errorf("clickhouse: unsupported column type: %s", field.Type.Name())
		}
	}

	for i := 0; i < itemsCount; i++ {
		for j := 0; j < columnCount; j++ {
			columnValue := value.Index(i).Field(columnIndexes[j]).Interface()
			switch columnTypes[j] {
			case "Time":
				actualValue, ok := columnValue.(time.Time)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues, ok := columns[j].([]time.Time)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues[i] = actualValue
				columns[j] = actualValues
			case "string":
				actualValue, ok := columnValue.(string)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues, ok := columns[j].([]string)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues[i] = actualValue
				columns[j] = actualValues
			case "int8":
				actualValue, ok := columnValue.(int8)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues, ok := columns[j].([]int8)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues[i] = actualValue
				columns[j] = actualValues
			case "int16":
				actualValue, ok := columnValue.(int16)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues, ok := columns[j].([]int16)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues[i] = actualValue
				columns[j] = actualValues
			case "int32":
				actualValue, ok := columnValue.(int32)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues, ok := columns[j].([]int32)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues[i] = actualValue
				columns[j] = actualValues
			case "int64", "int":
				actualValue, ok := columnValue.(int64)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues, ok := columns[j].([]int64)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues[i] = actualValue
				columns[j] = actualValues
			case "uint8":
				actualValue, ok := columnValue.(uint8)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues, ok := columns[j].([]uint8)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues[i] = actualValue
				columns[j] = actualValues
			case "uint16":
				actualValue, ok := columnValue.(uint16)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues, ok := columns[j].([]uint16)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues[i] = actualValue
				columns[j] = actualValues
			case "uint32":
				actualValue, ok := columnValue.(uint32)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues, ok := columns[j].([]uint32)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues[i] = actualValue
				columns[j] = actualValues
			case "uint64":
				actualValue, ok := columnValue.(uint64)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues, ok := columns[j].([]uint64)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues[i] = actualValue
				columns[j] = actualValues
			case "float32":
				actualValue, ok := columnValue.(float32)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues, ok := columns[j].([]float32)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues[i] = actualValue
				columns[j] = actualValues
			case "float64":
				actualValue, ok := columnValue.(float64)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues, ok := columns[j].([]float64)
				if !ok {
					return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
				}
				actualValues[i] = actualValue
				columns[j] = actualValues
			default:
				return fmt.Errorf("clickhouse: unsupported column type: %s", columnTypes[j])
			}
		}
	}

	batch, err := c.PrepareBatch(ctx, query)
	if err != nil {
		return err
	}

	for i := 0; i < columnCount; i++ {
		if err = batch.Column(i).Append(columns[i]); err != nil {
			return err
		}
	}

	return batch.Send()
}
