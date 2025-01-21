package notify

import (
	log "github.com/sirupsen/logrus"
)

type NotificationService struct {
	repo   Repo
	signal chan struct{}
}

func NewNotificationService(repo Repo, signal chan struct{}) *NotificationService {
	return &NotificationService{
		repo:   repo,
		signal: signal,
	}
}

func (n *NotificationService) Run() error {
	log.Debug(signalWait)
	select {
	case <-n.signal:
		log.Debug(signalRecieved)
		notify(n.repo)
	}
	return nil
}

func notify(repo Repo) {
	ul := UsersList{}
	list, err := ul.GetList(repo)
	if err != nil {
		log.WithField(errMsg, err).Error(readUserListError)
		return
	}

	if len(*list) < 1 {
		return
	}

	var userList []string
	response := make(map[string]Response)

	for _, user := range *list {
		userList = append(userList, user.UserUuid)

		if resp, ok := response[user.UserUuid]; ok {
			resp.Notifiers = append(resp.Notifiers, NotifierRecord{Target: user.Target, Origin: user.Origin})
			response[user.UserUuid] = resp
		} else {
			response[user.UserUuid] = Response{
				Notifiers: []NotifierRecord{{Target: user.Target, Origin: user.Origin}},
			}
		}
	}

	pl := PricesList{}
	list2, err := pl.GetList(repo, userList)
	if err != nil {
		log.WithField(errMsg, err).Error(readPriceListError)
		return
	}

	var siteMsgs []NotifyMessage

	for i := 0; i < len(*list2); i += 2 {
		if (*list2)[i].MerchUuid == (*list2)[i+1].MerchUuid {
			if (*list2)[i].Price != 0 && (*list2)[i+1].Price == 0 {
				// forming response for notification service
				if resp, ok := response[(*list2)[i].UserUuid]; ok {
					resp.Prices = append(resp.Prices, PriceRecords{
						MerchUuid: (*list2)[i].MerchUuid,
						Price:     (*list2)[i].Price,
					})
					response[(*list2)[i].UserUuid] = resp
				}

				// forming on site notifications
				siteMsgs = append(siteMsgs, NotifyMessage{
					UserUuid:  (*list2)[i].UserUuid,
					MerchUuid: (*list2)[i].MerchUuid,
					Price:     (*list2)[i].Price,
				})
			}
		}
	}

	err = CreateNotifications(repo, siteMsgs)
	if err != nil {
		log.WithField(errMsg, err).Error(createNotificationsError)
	}
}

func CreateNotifications(repo Repo, payload []NotifyMessage) error {
	err := repo.Create(payload)
	if err != nil {
		return err
	}
	return nil
}
