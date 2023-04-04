// Code generated by smithy-go-codegen DO NOT EDIT.

package sesv2

import (
	"context"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Returns the requested sending authorization policies for the given identity (an
// email address or a domain). The policies are returned as a map of policy names
// to policy contents. You can retrieve a maximum of 20 policies at a time. This
// API is for the identity owner only. If you have not verified the identity, this
// API will return an error. Sending authorization is a feature that enables an
// identity owner to authorize other senders to use its identities. For information
// about using sending authorization, see the Amazon SES Developer Guide
// (https://docs.aws.amazon.com/ses/latest/DeveloperGuide/sending-authorization.html).
// You can execute this operation no more than once per second.
func (c *Client) GetEmailIdentityPolicies(ctx context.Context, params *GetEmailIdentityPoliciesInput, optFns ...func(*Options)) (*GetEmailIdentityPoliciesOutput, error) {
	if params == nil {
		params = &GetEmailIdentityPoliciesInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "GetEmailIdentityPolicies", params, optFns, c.addOperationGetEmailIdentityPoliciesMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*GetEmailIdentityPoliciesOutput)
	out.ResultMetadata = metadata
	return out, nil
}

// A request to return the policies of an email identity.
type GetEmailIdentityPoliciesInput struct {

	// The email identity.
	//
	// This member is required.
	EmailIdentity *string

	noSmithyDocumentSerde
}

// Identity policies associated with email identity.
type GetEmailIdentityPoliciesOutput struct {

	// A map of policy names to policies.
	Policies map[string]string

	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationGetEmailIdentityPoliciesMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsRestjson1_serializeOpGetEmailIdentityPolicies{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestjson1_deserializeOpGetEmailIdentityPolicies{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = addOpGetEmailIdentityPoliciesValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opGetEmailIdentityPolicies(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addRequestIDRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	return nil
}

func newServiceMetadataMiddleware_opGetEmailIdentityPolicies(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "ses",
		OperationName: "GetEmailIdentityPolicies",
	}
}
