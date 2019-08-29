package marketplace

import (
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/internal/errors"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// getLocalImage returns the correct local version of an image matching
// the current zone and the compatible commercial type
func (version *Version) getLocalImage(zone scw.Zone, commercialType string) (*LocalImage, error) {

	for _, localImage := range version.LocalImages {

		// Check if in correct zone
		if localImage.Zone != zone {
			continue
		}

		// Check if compatible with wanted commercial type
		for _, compatibleCommercialType := range localImage.CompatibleCommercialTypes {
			if compatibleCommercialType == commercialType {
				return localImage, nil
			}
		}
	}

	return nil, fmt.Errorf("couldn't find compatible local image for this image version (%s)", version.ID)

}

// getLatestVersion returns the current/latests version on an image,
// or an error in case the image doesn't have a public version.
func (image *Image) getLatestVersion() (*Version, error) {

	for _, version := range image.Versions {
		if version.ID == image.CurrentPublicVersion {
			return version, nil
		}
	}

	return nil, errors.New("latest version could not be found for image %s", image.Name)
}

// GetLocalImageIDByNameRequest is used by FindLocalImageIDByName
type GetLocalImageIDByNameRequest struct {
	ImageLabel     string
	Zone           scw.Zone
	CommercialType string
}

// GetLocalImageIDByLabel search for an image with the given name (exact match) in the given region
// it returns the latest version of this specific image.
func (s *API) GetLocalImageIDByLabel(req *GetLocalImageIDByNameRequest) (string, error) {

	listImageRequest := &ListImagesRequest{}
	listImageResponse, err := s.ListImages(listImageRequest, scw.WithAllPages())
	if err != nil {
		return "", err
	}

	images := listImageResponse.Images
	label := strings.Replace(req.ImageLabel, "-", "_", -1)

	for _, image := range images {

		// Match name of the image
		if label == image.Label {

			latestVersion, err := image.getLatestVersion()
			if err != nil {
				return "", errors.Wrap(err, "couldn't find a matching image for the given name (%s), zone (%s) and commercial type (%s)", req.ImageLabel, req.Zone, req.CommercialType)
			}

			localImage, err := latestVersion.getLocalImage(req.Zone, req.CommercialType)
			if err != nil {
				return "", errors.Wrap(err, "couldn't find a matching image for the given name (%s), zone (%s) and commercial type (%s)", req.ImageLabel, req.Zone, req.CommercialType)
			}

			return localImage.ID, nil
		}

	}

	return "", errors.New("couldn't find a matching image for the given name (%s), zone (%s) and commercial type (%s)", req.ImageLabel, req.Zone, req.CommercialType)
}

// UnsafeSetTotalCount should not be used
// Internal usage only
func (r *ListImagesResponse) UnsafeSetTotalCount(totalCount int) {
	r.TotalCount = uint32(totalCount)
}

// UnsafeSetTotalCount should not be used
// Internal usage only
func (r *ListVersionsResponse) UnsafeSetTotalCount(totalCount int) {
	r.TotalCount = uint32(totalCount)
}
