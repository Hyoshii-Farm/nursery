package predator

import (
	"log"
	"math"

	model "github.com/Hyoshii-Farm/nursery/feature/report/predator/models"
)

func (s *Service) GetReport(q *model.PredatorPageQuery) (*model.PredatorPageResponse, error) {
	// Dates are already parsed and validated in DTO
	start := q.StartDate
	end := q.EndDate

	duration := end.Sub(start)
	prevEnd := start.AddDate(0, 0, -1)
	prevStart := prevEnd.Add(-duration)

	// Fetch current data
	dateLayout := "2006-01-02"
	history, err := s.repo.GetHistory(q.StartDate.Format(dateLayout), q.EndDate.Format(dateLayout), q.PredatorIDs, q.PageHistory, q.Limit)
	if err != nil {
		log.Printf("[predator] GetHistory error: %v", err)
	}
	historyTotal, err := s.repo.CountHistory(q.StartDate.Format(dateLayout), q.EndDate.Format(dateLayout), q.PredatorIDs)
	if err != nil {
		log.Printf("[predator] CountHistory error: %v", err)
	}
	currVariants, err := s.repo.GetVariantSummary(q.EndDate.Format(dateLayout), q.PredatorIDs)
	if err != nil {
		log.Printf("[predator] GetVariantSummary (current) error: %v", err)
	}
	currActions, err := s.repo.GetActionSummary(q.StartDate.Format(dateLayout), q.EndDate.Format(dateLayout), q.PredatorIDs)
	if err != nil {
		log.Printf("[predator] GetActionSummary (current) error: %v", err)
	}

	// Fetch previous period data for comparison
	prevDate := prevEnd.Format("2006-01-02")
	prevVariants, err := s.repo.GetVariantSummary(prevDate, q.PredatorIDs)
	if err != nil {
		log.Printf("[predator] GetVariantSummary (prev) error: %v", err)
	}
	prevActions, err := s.repo.GetActionSummary(prevStart.Format("2006-01-02"), prevDate, q.PredatorIDs)
	if err != nil {
		log.Printf("[predator] GetActionSummary (prev) error: %v", err)
	}

	// Helper function for gap calculation
	calcGap := func(curr, prev float64) float64 {
		if prev == 0 {
			if curr > 0 {
				return 100
			}
			return 0
		}
		return ((curr - prev) / prev) * 100
	}

	// 1. Total Stock KPI
	var currTotal model.KPIBlock
	var currTotalQty float64
	for name, v := range currVariants {
		currTotalQty += v.Quantity
		if v.Quantity > currTotal.HighestQuantity {
			currTotal.HighestQuantity = v.Quantity
			currTotal.HighestName = name
		}
		if currTotal.LowestName == "" || v.Quantity < currTotal.LowestQuantity {
			currTotal.LowestQuantity = v.Quantity
			currTotal.LowestName = name
		}
	}
	var prevTotalQty float64
	for _, v := range prevVariants {
		prevTotalQty += v.Quantity
	}
	currTotal.CurrentQuantity = currTotalQty
	currTotal.LastPeriodQuantity = prevTotalQty
	currTotal.GapPercentage = calcGap(currTotalQty, prevTotalQty)

	// 2. Action KPIs (New/Removed)
	newStock := currActions["ADD"]
	prevNew := prevActions["ADD"]
	newStock.LastPeriodQuantity = prevNew.CurrentQuantity
	newStock.GapPercentage = calcGap(newStock.CurrentQuantity, prevNew.CurrentQuantity)

	removedStock := currActions["REMOVED"]
	prevRemoved := prevActions["REMOVED"]
	removedStock.LastPeriodQuantity = prevRemoved.CurrentQuantity
	removedStock.GapPercentage = calcGap(removedStock.CurrentQuantity, prevRemoved.CurrentQuantity)

	// Helper function for rounding to 1 decimal place
	round := func(val float64) float64 {
		return math.Round(val*10) / 10
	}

	// 3. Average Age KPI
	var currAge model.AgeBlock
	var totalAge, totalPrevAge float64
	if len(currVariants) > 0 {
		currAge.YoungestName = "None" // Youngest initial
		for name, v := range currVariants {
			totalAge += v.Oldest
			if v.Oldest > currAge.Oldest {
				currAge.Oldest = v.Oldest
				currAge.OldestName = name
			}
			if currAge.YoungestName == "None" || v.Youngest < currAge.Youngest {
				currAge.Youngest = v.Youngest
				currAge.YoungestName = name
			}
		}
		currAge.CurrentAge = round(totalAge / float64(len(currVariants)))
		currAge.Oldest = round(currAge.Oldest)
		currAge.Youngest = round(currAge.Youngest)
	}
	if len(prevVariants) > 0 {
		for _, v := range prevVariants {
			totalPrevAge += v.Oldest
		}
		currAge.LastPeriodAge = round(totalPrevAge / float64(len(prevVariants)))
		currAge.GapPercentage = calcGap(currAge.CurrentAge, currAge.LastPeriodAge)
	}

	// 4. Variant KPIs
	variantKPIs := []model.VariantKPI{}
	for name, currV := range currVariants {
		prevV := prevVariants[name]
		variantKPIs = append(variantKPIs, model.VariantKPI{
			VariantName:        name,
			Quantity:           currV.Quantity,
			LastPeriodQuantity: prevV.Quantity,
			GapQuantity:        currV.Quantity - prevV.Quantity,
			Oldest:             round(currV.Oldest),
			Youngest:           round(currV.Youngest),
		})
	}

	// 5. Stock Section Data (paginated)
	stockItems, stockTotal, err := s.repo.GetStockPaginated(q.EndDate.Format(dateLayout), q.PredatorIDs, q.PageStock, q.Limit)
	if err != nil {
		log.Printf("[predator] GetStockPaginated error: %v", err)
	}
	// Round stock ages
	for i := range stockItems {
		stockItems[i].Age = round(stockItems[i].Age)
	}

	// 6. History mapping
	historyItems := []model.HistoryItem{}
	for _, h := range history {
		p, _ := s.repo.FindPredator(h.PredatorID)
		historyItems = append(historyItems, model.HistoryItem{
			Date:        h.Datetime,
			Action:      h.Action,
			VariantName: p.Name,
			Quantity:    h.Quantity,
			Note:        h.Pic,
		})
	}

	response := &model.PredatorPageResponse{
		KPI: model.KPISection{
			TotalStock:   currTotal,
			NewStock:     newStock,
			RemovedStock: removedStock,
			AverageAge:   currAge,
			Variant:      variantKPIs,
		},
		Stock: model.StockSection{
			Pagination: calcPagination(int(stockTotal), q.PageStock, q.Limit),
			Data:       stockItems,
		},
		History: model.HistorySection{
			Pagination: calcPagination(int(historyTotal), q.PageHistory, q.Limit),
			Data:       historyItems,
		},
	}

	return response, nil
}
