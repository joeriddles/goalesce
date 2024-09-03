// TODO: refactor all of this V1 code
package api

import (
	"bufio"
	"bytes"
	"fmt"
	"go/types"
	"strings"

	"github.com/joeriddles/goalesce/pkg/v2/entity"
	"github.com/joeriddles/goalesce/pkg/v2/utils"
)

// Convert the field be converted to matching field on dst
func convertField(
	field *entity.GormModelField,
	dst *entity.GormModelMetadata,
) string {
	from := "src"
	to := "dst"

	dstField := dst.GetField(field.Name)

	if field.MapApiFunc != nil {
		return fmt.Sprintf("%v.%v = model.%v(%v.%v)", to, dstField.Name, *field.MapApiFunc, from, field.Name)
	}

	srcType := field.GetType()
	dstType := dstField.GetType()

	isSrcPtr := false
	if ptrSrc, ok := srcType.(*types.Pointer); ok {
		isSrcPtr = true
		srcType = ptrSrc.Elem()
	}

	isDstPtr := false
	if ptrDst, ok := dstType.(*types.Pointer); ok {
		isDstPtr = true
		dstType = ptrDst.Elem()
	}

	switch s := srcType.(type) {
	case *types.Basic:
		switch d := dstType.(type) {
		case *types.Basic:
			if s.Kind() != d.Kind() && types.ConvertibleTo(s, d) {
				if isSrcPtr && isDstPtr {
					var b bytes.Buffer
					w := bufio.NewWriter(&b)
					w.WriteString(fmt.Sprintf("dst.%v = func() *%v {", dstField.Name, d.Name()))
					w.WriteString(fmt.Sprintf("	var %v %v", utils.ToCamelCase(dstField.Name), d.Name()))
					w.WriteString(fmt.Sprintf("	if src.%v != nil {", field.Name))
					w.WriteString(fmt.Sprintf("		%v = %v(*src.%v)", utils.ToCamelCase(dstField.Name), d.Name(), field.Name))
					w.WriteString(fmt.Sprintf("	}%v", "")) // empty string to keep linter happy
					w.WriteString(fmt.Sprintf("	return &%v", utils.ToCamelCase(dstField.Name)))
					w.WriteString(fmt.Sprintf("}()%v", ""))
					if err := w.Flush(); err != nil {
						return err.Error()
					}
					return b.String()
				}

				return fmt.Sprintf("%v.%v = %v(%v.%v)", to, dstField.Name, d.Name(), from, field.Name)
			}
		}
	case *types.Named:
		switch d := dstType.(type) {
		case *types.Named:
			if s.Obj().Name() == "Time" && d.Obj().Name() == "DeletedAt" {
				return fmt.Sprintf("%v.%v = convertTimeToGormDeletedAt(%v.%v)", to, dstField.Name, from, field.Name)
			} else if d.Obj().Name() == "Time" && s.Obj().Name() == "DeletedAt" {
				return fmt.Sprintf("%v.%v = convertGormDeletedAtToTime(%v.%v)", to, dstField.Name, from, field.Name)
			}

			// TODO(joeriddles): add field to GormModelField for references to user-defined models?
			if utils.IsComplexType(dstField.Type) && !strings.Contains(dstField.Type, ".") {
				isSrcPtr := strings.Contains(field.Type, "*")
				mapperName, isDstPtr := strings.CutPrefix(dstField.Type, "*")
				if dst.IsApi {
					mapperName = mapperName + "Api"
				}

				if isDstPtr {
					if !isSrcPtr {
						from = "&" + from
					}
					return fmt.Sprintf(`%v.%v = New%vMapper().MapPtr(%v.%v)`, to, dstField.Name, mapperName, from, field.Name)
				} else {
					return fmt.Sprintf(`%v.%v = New%vMapper().Map(%v.%v)`, to, dstField.Name, mapperName, from, field.Name)
				}
			}
		}
	case *types.Slice:
		if _, ok := dstType.(*types.Slice); ok {
			isDstPtr := strings.HasPrefix(dstField.Type, "*")
			var isDstElemPtr bool
			if isDstPtr {
				isDstElemPtr = dstField.Type[3:4] == "*"
			} else {
				isDstElemPtr = dstField.Type[2:3] == "*"
			}

			isSrcPtr := strings.HasPrefix(field.Type, "*")
			var isSrcElemPtr bool
			if isSrcPtr {
				isSrcElemPtr = field.Type[3:4] == "*"
			} else {
				isSrcElemPtr = field.Type[2:3] == "*"
			}

			mapperName := strings.ReplaceAll(strings.ReplaceAll(dstField.Type, "*", ""), "[]", "")
			if dst.IsApi {
				mapperName = mapperName + "Api"
			}

			if dst.IsApi {
				if isSrcPtr && isSrcElemPtr {
					return fmt.Sprintf(`if %v.%v != nil { %v.%v = New%vMapper().MapPtrSlicePtrs(%v.%v) }`, from, field.Name, to, dstField.Name, mapperName, from, field.Name)
				} else if isSrcPtr {
					return fmt.Sprintf(`if %v.%v != nil { %v.%v = New%vMapper().MapPtrSlice(%v.%v) }`, from, field.Name, to, dstField.Name, mapperName, from, field.Name)
				} else if isSrcElemPtr {
					return fmt.Sprintf(`if %v.%v != nil { %v.%v = New%vMapper().MapSlicePtrs(%v.%v) }`, from, field.Name, to, dstField.Name, mapperName, from, field.Name)
				} else {
					return fmt.Sprintf(`if %v.%v != nil { %v.%v = New%vMapper().MapSlice(%v.%v) }`, from, field.Name, to, dstField.Name, mapperName, from, field.Name)
				}
			} else {
				if isDstPtr && isDstElemPtr {
					return fmt.Sprintf(`if %v.%v != nil { %v.%v = New%vMapper().MapPtrSlicePtrs(%v.%v) }`, from, field.Name, to, dstField.Name, mapperName, from, field.Name)
				} else if isDstPtr {
					return fmt.Sprintf(`if %v.%v != nil { %v.%v = New%vMapper().MapPtrSlice(%v.%v) }`, from, field.Name, to, dstField.Name, mapperName, from, field.Name)
				} else if isDstElemPtr {
					return fmt.Sprintf(`if %v.%v != nil { %v.%v = New%vMapper().MapSlicePtrs(%v.%v) }`, from, field.Name, to, dstField.Name, mapperName, from, field.Name)
				} else {
					return fmt.Sprintf(`if %v.%v != nil { %v.%v = New%vMapper().MapSlice(%v.%v) }`, from, field.Name, to, dstField.Name, mapperName, from, field.Name)
				}
			}
		}
	}

	return fmt.Sprintf("%v.%v = %v.%v", to, dstField.Name, from, field.Name)
}
