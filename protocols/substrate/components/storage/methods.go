package substrate

import (
	"context"
	"fmt"
	"io"

	"github.com/fxamacker/cbor/v2"
	"github.com/ipfs/go-cid"
	hoarderIface "github.com/taubyte/go-interfaces/services/hoarder"
	storageIface "github.com/taubyte/go-interfaces/services/substrate/components/storage"
	hoarderSpecs "github.com/taubyte/go-specs/hoarder"
	"github.com/taubyte/p2p/peer"
	"github.com/taubyte/tau/protocols/substrate/components/storage/common"
)

// Context Config will be generated by New()
func (s *Service) Storage(context storageIface.Context) (storageIface.Storage, error) {
	hash, err := common.GetStorageHash(context)
	if err != nil {
		return nil, err
	}

	s.storagesLock.RLock()
	storage, ok := s.storages[hash]
	s.storagesLock.RLocker().Unlock()
	if !ok {
		context.Config, err = s.getStoreConfig(context.ProjectId, context.ApplicationId, context.Matcher)
		if err != nil {
			return nil, err
		}

		storage, err = s.storageMethod(s, s.dbFactory, context, common.Logger, s.matcher)
		if err != nil {
			return nil, err
		}

		s.storagesLock.Lock()
		s.storages[hash] = storage
		s.storagesLock.Unlock()

		err = s.pubsubStorage(context, s.Branch())
		if err != nil {
			return nil, fmt.Errorf("pubsub storage `%s` failed with: %s", context.Matcher, err)
		}

		commit, err := s.Tns().Simple().Commit(context.ProjectId, s.Branch())
		if err != nil {
			return nil, fmt.Errorf("getting commit for project id `%s` and branch `%s` failed with: %s", context.ProjectId, s.Branch(), err)
		}

		s.commitLock.Lock()
		s.commits[hash] = commit
		s.commitLock.Unlock()
	}

	valid, newCommitId, err := s.validateCommit(hash, context.ProjectId, s.Branch())
	if err != nil {
		return nil, err
	}

	if !valid {
		s.storagesLock.Lock()
		s.commitLock.Lock()

		defer s.storagesLock.Unlock()
		defer s.commitLock.Unlock()

		storage, err = s.updateStorage(storage)
		if err != nil {
			return nil, err
		}

		s.storages[hash] = storage
		s.commits[hash] = newCommitId
	}

	return storage, nil
}

func (s *Service) Add(content io.Reader) (cid.Cid, error) {
	__cid, err := s.Node().AddFile(content)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("failed adding file with %v", err)
	}

	_cid, err := cid.Parse(__cid)
	if err != nil {
		return cid.Cid{}, fmt.Errorf("failed parsing cid with %v", err)
	}

	return _cid, nil
}

func (s *Service) GetFile(ctx context.Context, cid cid.Cid) (peer.ReadSeekCloser, error) {
	file, err := s.Node().GetFileFromCid(ctx, cid)
	if err != nil {
		return nil, fmt.Errorf("failed grabbing file %s from ipfs with %w", cid, err)
	}

	return file, nil
}
func (s *Service) pubsubStorage(context storageIface.Context, branch string) error {
	auction := &hoarderIface.Auction{
		Type:     hoarderIface.AuctionNew,
		MetaType: hoarderIface.Storage,
		Meta: hoarderIface.MetaData{
			ConfigId:      context.Config.Id,
			ApplicationId: context.ApplicationId,
			ProjectId:     context.ProjectId,
			Match:         context.Matcher,
			Branch:        s.Branch(),
		},
	}

	dataBytes, err := cbor.Marshal(auction)
	if err != nil {
		return fmt.Errorf("marshalling auction failed with %w", err)
	}

	topic, err := s.Node().Messaging().Join(hoarderSpecs.PubSubIdent)
	if err != nil {
		return fmt.Errorf("getting topic for `%s` failed with: %w", hoarderSpecs.PubSubIdent, err)
	}

	if err = topic.Publish(s.Context(), dataBytes); err != nil {
		return fmt.Errorf("failed publishing storage %s with %w", context.Matcher, err)
	}

	return nil
}
