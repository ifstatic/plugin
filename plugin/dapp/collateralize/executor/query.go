// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/types"
	pty "github.com/33cn/plugin/plugin/dapp/collateralize/types"
)

func (c *Collateralize) Query_CollateralizeInfoByID(req *pty.ReqCollateralizeInfo) (types.Message, error) {
	coll,err := queryCollateralizeByID(c.GetStateDB(), req.CollateralizeId)
	if err != nil {
		clog.Error("Query_CollateralizeInfoByID", "id", req.CollateralizeId, "error", err)
		return nil, err
	}

	info :=  &pty.RepCollateralizeCurrentInfo{
		Status:             coll.Status,
		TotalBalance:       coll.TotalBalance,
		DebtCeiling:        coll.DebtCeiling,
		LiquidationRatio:   coll.LiquidationRatio,
		StabilityFeeRatio:  coll.StabilityFeeRatio,
		CreateAddr:         coll.CreateAddr,
		Balance:            coll.Balance,
		Period:             coll.Period,
		CollateralizeId:    coll.CollateralizeId,
		CollBalance:        coll.CollBalance,
	}
	info.BorrowRecords = append(info.BorrowRecords, coll.BorrowRecords...)
	info.BorrowRecords = append(info.BorrowRecords, coll.InvalidRecords...)

	return info, nil
}

func (c *Collateralize) Query_CollateralizeInfoByIDs(req *pty.ReqCollateralizeInfos) (types.Message, error) {
	infos := &pty.RepCollateralizeCurrentInfos{}
	for _, id := range req.CollateralizeIds {
		coll,err := queryCollateralizeByID(c.GetStateDB(), id)
		if err != nil {
			clog.Error("Query_CollateralizeInfoByID", "id", id, "error", err)
			return nil, err
		}

		info :=  &pty.RepCollateralizeCurrentInfo{
			Status:             coll.Status,
			TotalBalance:       coll.TotalBalance,
			DebtCeiling:        coll.DebtCeiling,
			LiquidationRatio:   coll.LiquidationRatio,
			StabilityFeeRatio:  coll.StabilityFeeRatio,
			CreateAddr:         coll.CreateAddr,
			Balance:            coll.Balance,
			Period:             coll.Period,
			CollateralizeId:    coll.CollateralizeId,
			CollBalance:        coll.CollBalance,
		}
		info.BorrowRecords = append(info.BorrowRecords, coll.BorrowRecords...)
		info.BorrowRecords = append(info.BorrowRecords, coll.InvalidRecords...)

		infos.Infos = append(infos.Infos, info)
	}

	return infos, nil
}

func (c *Collateralize) Query_CollateralizeByStatus(req *pty.ReqCollateralizeByStatus) (types.Message, error) {
	ids := &pty.RepCollateralizeIDs{}
	collIDRecords, err := queryCollateralizeByStatus(c.GetLocalDB(), req.Status, req.Index)
	if err != nil {
		clog.Error("Query_CollateralizeByStatus", "get collateralize record error", err)
		return nil, err
	}

	ids.IDs = append(ids.IDs, collIDRecords...)
	return ids, nil
}

func (c *Collateralize) Query_CollateralizeByAddr(req *pty.ReqCollateralizeByAddr) (types.Message, error) {
	ids := &pty.RepCollateralizeIDs{}
	collIDRecords, err := queryCollateralizeByAddr(c.GetLocalDB(), req.Addr, req.Status, req.Index)
	if err != nil {
		clog.Error("Query_CollateralizeByAddr", "get collateralize record error", err)
		return nil, err
	}

	ids.IDs = append(ids.IDs, collIDRecords...)
	return ids, nil
}

func (c *Collateralize) Query_CollateralizeRecordByID(req *pty.ReqCollateralizeRecord) (types.Message, error) {
	ret := &pty.RepCollateralizeRecord{}
	issuRecord, err := queryCollateralizeRecordByID(c.GetStateDB(), req.CollateralizeId, req.RecordId)
	if err != nil {
		clog.Error("Query_IssuanceRecordByID", "get collateralize record error", err)
		return nil, err
	}

	ret.Record = issuRecord
	return ret, nil
}

func (c *Collateralize) Query_CollateralizeRecordByAddr(req *pty.ReqCollateralizeRecordByAddr) (types.Message, error) {
	ret := &pty.RepCollateralizeRecords{}
	records, err := queryCollateralizeRecordByAddr(c.GetStateDB(), c.GetLocalDB(), req.Addr, req.Index)
	if err != nil {
		clog.Error("Query_CollateralizeRecordByAddr", "get collateralize record error", err)
		return nil, err
	}

	if req.Status == 0 {
		ret.Records = records
	} else {
		for _,record := range records {
			if record.Status == req.Status {
				ret.Records = append(ret.Records, record)
			}
		}
	}
	return ret, nil
}

func (c *Collateralize) Query_CollateralizeRecordByStatus(req *pty.ReqCollateralizeRecordByStatus) (types.Message, error) {
	ret := &pty.RepCollateralizeRecords{}
	records, err := queryCollateralizeRecordByStatus(c.GetStateDB(), c.GetLocalDB(), req.Status, req.Index)
	if err != nil {
		clog.Error("Query_CollateralizeRecordByStatus", "get collateralize record error", err)
		return nil, err
	}

	ret.Records = records
	return ret, nil
}

func (c *Collateralize) Query_CollateralizeConfig(req *pty.ReqCollateralizeRecordByAddr) (types.Message, error) {
	config, err := getCollateralizeConfig(c.GetStateDB())
	if err != nil {
		clog.Error("Query_CollateralizeConfig", "get collateralize config error", err)
		return nil, err
	}

	collIDRecords, err := queryCollateralizeByStatus(c.GetLocalDB(), pty.CollateralizeStatusCreated, 0)
	if err != nil {
		clog.Debug("Query_CollateralizeByStatus", "get collateralize record error", err)
	}

	balance := config.TotalBalance
	for _, id := range collIDRecords {
		coll, err := queryCollateralizeByID(c.GetStateDB(), id)
		if err != nil {
			clog.Error("Query_CollateralizeInfoByID", "id", id, "error", err)
			return nil, err
		}

		balance -= coll.TotalBalance
	}
	ret := &pty.RepCollateralizeConfig{
		TotalBalance:config.TotalBalance,
		DebtCeiling: config.DebtCeiling,
		LiquidationRatio: config.LiquidationRatio,
		StabilityFeeRatio: config.StabilityFeeRatio,
		Period: config.Period,
		Balance: balance,
		CurrentTime: config.CurrentTime,
	}

	return ret, nil
}

func (c *Collateralize) Query_CollateralizePrice(req *pty.ReqCollateralizeRecordByAddr) (types.Message, error) {
	price, err := getLatestPrice(c.GetStateDB())
	if err != nil {
		clog.Error("Query_CollateralizePrice", "error", err)
		return nil, err
	}

	return &pty.RepCollateralizePrice{Price:price}, nil
}