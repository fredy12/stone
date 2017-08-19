package quota

const (
	QuotaBare = "bare"
	QuotaXfs  = "xfs"
	QuotaExt4 = "ext4"
)

type Quota struct {
	Size uint64
}

type newTypeQuota func(basePath string) (QuotaControl, error)

type QuotaControl interface {
	Name() string
	SetQuota(targetPath string, quota *Quota) error
	GetQuota(targetPath string) (*Quota, error)
	RemoveQuota(targetPath string) error
}

var (
	quotaControls    map[string]QuotaControl
	registerNewQuota map[string]newTypeQuota
)

func init() {
	registerNewQuota = map[string]newTypeQuota{
		QuotaBare: NewBareQuotaControl,
		QuotaXfs:  NewXfsQuotaControl,
		QuotaExt4: NewExt4QuotaControl,
	}
	quotaControls = map[string]QuotaControl{}
}

// NewQuota return quotaControl. basePath is the path to store all volumes
func NewQuota(basePath string, format string) (QuotaControl, error) {
	if v, exists := quotaControls[basePath]; exists {
		if v.Name() == format {
			// exist
			return v, nil
		}
		// else exist but wrong format, recreate it
	}
	// recreate
	if newQuota, exist := registerNewQuota[format]; exist {
		c, err := newQuota(basePath)
		if err != nil {
			return nil, err
		}
		quotaControls[basePath] = c
	} else {
		// use bare
		c, _ := registerNewQuota[QuotaBare](basePath)
		quotaControls[basePath] = c
	}
	return quotaControls[basePath], nil
}
