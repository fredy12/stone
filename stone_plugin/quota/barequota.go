package quota

type BareQuotaControl struct{}

func NewBareQuotaControl(basePath string) (QuotaControl, error) {
	return &BareQuotaControl{}, nil
}

func (q *BareQuotaControl) Name() string {
	return QuotaBare
}

func (q *BareQuotaControl) SetQuota(targetPath string, quota *Quota) error {
	return nil
}

func (q *BareQuotaControl) GetQuota(targetPath string) (*Quota, error) {
	return &Quota{Size: 0}, nil
}

func (q *BareQuotaControl) RemoveQuota(targetPath string) error {
	return nil
}
