package handler

import (
	"digicon/common/errors"
	proto "digicon/proto/rpc"
	"digicon/public_service/model"
	"encoding/json"
	"golang.org/x/net/context"
)

func (s *RPCServer) GetTokensList(ctx context.Context, req *proto.TokensRequest, rsp *proto.TokensResponse) error {
	u := model.Tokens{}
	r := u.GetTokens(req.Tokens)
	for _, v := range r {

		rsp.Tokens = append(rsp.Tokens, &proto.TokensData{
			TokenId: int32(v.Id),
			Mark:    v.Mark,
		})
	}
	return nil
}

func (s *RPCServer) GetSiteConfig(ctx context.Context, req *proto.NullRequest, rsp *proto.GetSiteConfigResponse) error {
	configModel := new(model.ConfigModel)

	// 1.基础配置
	config, err := configModel.Get(model.SITE_CONFIG_NAME_SITE)
	if err != nil {
		rsp.Code = int32(errors.GetErrStatus(err))
		rsp.Message = errors.GetErrMsg(err)
		return nil
	}
	// 解析
	siteConfig := new(model.SiteConfig)
	err = json.Unmarshal([]byte(config.Value), siteConfig)
	if err != nil {
		err = errors.NewSys(err)

		rsp.Code = int32(errors.GetErrStatus(err))
		rsp.Message = errors.GetErrMsg(err)
		return nil
	}

	// 2.客服配置
	config, err = configModel.Get(model.SITE_CONFIG_NAME_KEFU)
	if err != nil {
		rsp.Code = int32(errors.GetErrStatus(err))
		rsp.Message = errors.GetErrMsg(err)
		return nil
	}
	// 解析
	kefuConfig := new(model.KefuConfig)
	err = json.Unmarshal([]byte(config.Value), kefuConfig)
	if err != nil {
		err = errors.NewSys(err)

		rsp.Code = int32(errors.GetErrStatus(err))
		rsp.Message = errors.GetErrMsg(err)
		return nil
	}

	// 设置
	rsp.Data = &proto.GetSiteConfigResponse_Data{}
	rsp.Data.Site = &proto.GetSiteConfigResponse_Data_Site{
		Name:            siteConfig.Name,
		EnglishName:     siteConfig.EnglishName,
		Title:           siteConfig.Title,
		EnglishTitle:    siteConfig.EnglishTitle,
		Logo:            siteConfig.Logo,
		Keyword:         siteConfig.Keyword,
		EnglishKeyword:  siteConfig.EnglishKeyword,
		Desc:            siteConfig.Desc,
		EnglishDesc:     siteConfig.EnglishDesc,
		Beian:           siteConfig.Beian,
		StatisticScript: siteConfig.StatisticScript,
	}
	rsp.Data.Kefu = &proto.GetSiteConfigResponse_Data_Kefu{
		Phone:   kefuConfig.Phone,
		Email:   kefuConfig.Email,
		Address: kefuConfig.Address,
		Dianbao: kefuConfig.Dianbao,
	}

	return nil
}
