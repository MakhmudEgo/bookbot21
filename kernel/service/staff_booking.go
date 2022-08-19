package service

import (
	"context"

	"boobot/kernel/domain/models"
	"boobot/kernel/service/chainer"
	staffBkg "boobot/kernel/service/chainer/staff_booking"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type staffBooking struct {
	opts  *Opts
	chain chainer.Chainer
}

func NewStaffBooking(opts *Opts) (Service, error) {
	//todo: check opts
	return &staffBooking{
		opts: opts,
	}, nil
}

func (s staffBooking) Execute(ctx context.Context, user *models.User) (*tg.MessageConfig, error) {
	user.HandleStep = chainer.CheckStepHandle(user.HandleStep, chainer.StaffShowBtnBookingsStep,
		chainer.StaffBookingSteps...)

	if s.opts.Update.Message != nil && s.opts.Update.Message.Text == string(models.Staff) {
		user.HandleStep = int(chainer.StaffShowBtnBookingsStep)
	}

	opts := &chainer.Opts{
		UserRepo:    s.opts.UserRepo,
		Update:      s.opts.Update,
		SessionRepo: s.opts.SessionRepo,
		RootRepo:    s.opts.RootRepo,
	}

	chain := staffBkg.NewShowBtn(opts)
	//chain.SetNext(register.NewSendConfirmURL(opts))

	msgReply, err := chain.Handle(ctx, user)
	if err != nil {
		return nil, err
	}

	if s.opts.Update.Message != nil {
		msgReply.ChatID = s.opts.Update.Message.From.ID
	} else {
		msgReply.ChatID = s.opts.Update.CallbackQuery.From.ID
	}

	return msgReply, nil
}