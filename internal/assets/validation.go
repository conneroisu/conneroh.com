package assets

import (
	"strings"

	"github.com/rotisserie/eris"
)

// Defaults sets the default values for the document if they are missing.
func Defaults(doc *Doc) error {
	// Set default icon if not provided
	if doc.Icon == "" {
		doc.Icon = "tag"
	}

	return nil
}

// Validate validate the given embedding.
func Validate(
	path string,
	emb *Doc,
) error {
	var (
		err  error
		errs = []error{}
	)
	if emb.Slug == "" {
		errs = append(errs, eris.Wrapf(
			ErrValueMissing,
			"%s is missing slug",
			path,
		))
	}

	if emb.Description == "" {
		errs = append(errs, eris.Wrapf(
			ErrValueMissing,
			"%s is missing description",
			path,
		))
	}

	if emb.Content == "" {
		errs = append(errs, eris.Wrapf(
			ErrValueMissing,
			"%s is missing content",
			path,
		))
	}

	if strings.Contains(emb.Slug, " ") {
		errs = append(errs, eris.Wrapf(
			ErrValueInvalid,
			"slug %s contains spaces",
			path,
		))
	}

	// TODO: Validate that the banner path is a valid image and is not null/zero

	for _, er := range errs {
		if er != nil {
			err = eris.Wrapf(err, "failed validating %s", path)
		}
	}
	if err != nil {
		return err
	}

	return nil
}
