package ucloudstack

import (
	"context"
	"errors"
	"io"

	v1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog/v2"
)

const providerName = "ucloudstack"

var ErrCloudInstanceNotFound = errors.New("ucloudstack instance not found")

// 注册
func init() {
	cloudprovider.RegisterCloudProvider(providerName, NewCloudProvider)
}

func NewCloudProvider(config io.Reader) (cloudprovider.Interface, error) {
	return &UCloudStack{}, nil
}

type UCloudStack struct {
}

type Config struct {
	KubeConfig            string  `yaml:"kube_config" json:"kube_config"`
	Master                string  `yaml:"master" json:"master"`
	Qps                   float32 `yaml:"qps" json:"qps"`
	Burst                 int     `yaml:"burst" json:"burst"`
	InsecureSkipTlsVerify bool    `yaml:"insecure_skip_tls_verify" json:"insecure_skip_tls_verify"`
}

// https://github.com/kubernetes/cloud-provider/blob/master/cloud.go
// Initialize provides the cloud with a kubernetes client builder and may spawn goroutines
// to perform housekeeping or run custom controllers specific to the cloud provider.
// Any tasks started here should be cleaned up when the stop channel closes.
func (us *UCloudStack) Initialize(clientBuilder cloudprovider.ControllerClientBuilder, stop <-chan struct{}) {
	// clientset := clientBuilder.ClientOrDie("cloud-controller-manager")
	// us.kubeClient = clientset
	klog.Info("Initialize")
}

// LoadBalancer returns a balancer interface. Also returns true if the interface is supported, false otherwise.
func (us *UCloudStack) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	klog.Info("LoadBalancer")
	return &LBClass{}, true
}

// Instances returns an instances interface. Also returns true if the interface is supported, false otherwise.
func (us *UCloudStack) Instances() (cloudprovider.Instances, bool) {
	klog.Info("Instances")
	return nil, false
}

// InstancesV2 is an implementation for instances and should only be implemented by external cloud providers.
// Implementing InstancesV2 is behaviorally identical to Instances but is optimized to significantly reduce
// API calls to the cloud provider when registering and syncing nodes. Implementation of this interface will
// disable calls to the Zones interface. Also returns true if the interface is supported, false otherwise.
func (us *UCloudStack) InstancesV2() (cloudprovider.InstancesV2, bool) {
	klog.Info("InstancesV2")
	return nil, false
}

// Zones returns a zones interface. Also returns true if the interface is supported, false otherwise.
// DEPRECATED: Zones is deprecated in favor of retrieving zone/region information from InstancesV2.
// This interface will not be called if InstancesV2 is enabled.
func (us *UCloudStack) Zones() (cloudprovider.Zones, bool) {
	klog.Info("Zones")
	return nil, false
}

// Clusters returns a clusters interface.  Also returns true if the interface is supported, false otherwise.
func (us *UCloudStack) Clusters() (cloudprovider.Clusters, bool) {
	klog.Info("Clusters")
	return nil, false
}

// Routes returns a routes interface along with whether the interface is supported.
func (us *UCloudStack) Routes() (cloudprovider.Routes, bool) {
	klog.Info("Routes")
	return nil, false
}

// ProviderName returns the cloud provider ID.
func (us *UCloudStack) ProviderName() string {
	klog.Info("ProviderName")
	return providerName
}

// HasClusterID returns true if a ClusterID is required and set
func (us *UCloudStack) HasClusterID() bool {
	klog.Info("HasClusterID")
	return true
}

// LoadBalancer is used for creating and maintaining load balancers
type LoadBalancer struct {
}

// LoadBalancerOpts have the options to talk to Neutron LBaaSV2 or Octavia
type LoadBalancerOpts struct {
	Enabled               bool                `gcfg:"enabled"`              // if false, disables the controller
	LBVersion             string              `gcfg:"lb-version"`           // overrides autodetection. Only support v2.
	SubnetID              string              `gcfg:"subnet-id"`            // overrides autodetection.
	MemberSubnetID        string              `gcfg:"member-subnet-id"`     // overrides autodetection.
	NetworkID             string              `gcfg:"network-id"`           // If specified, will create virtual ip from a subnet in network which has available IP addresses
	FloatingNetworkID     string              `gcfg:"floating-network-id"`  // If specified, will create floating ip for loadbalancer, or do not create floating ip.
	FloatingSubnetID      string              `gcfg:"floating-subnet-id"`   // If specified, will create floating ip for loadbalancer in this particular floating pool subnetwork.
	FloatingSubnet        string              `gcfg:"floating-subnet"`      // If specified, will create floating ip for loadbalancer in one of the matching floating pool subnetworks.
	FloatingSubnetTags    string              `gcfg:"floating-subnet-tags"` // If specified, will create floating ip for loadbalancer in one of the matching floating pool subnetworks.
	LBClasses             map[string]*LBClass // Predefined named Floating networks and subnets
	LBMethod              string              `gcfg:"lb-method"` // default to ROUND_ROBIN.
	LBProvider            string              `gcfg:"lb-provider"`
	CreateMonitor         bool                `gcfg:"create-monitor"`
	MonitorMaxRetries     uint                `gcfg:"monitor-max-retries"`
	ManageSecurityGroups  bool                `gcfg:"manage-security-groups"`
	InternalLB            bool                `gcfg:"internal-lb"` // default false
	CascadeDelete         bool                `gcfg:"cascade-delete"`
	FlavorID              string              `gcfg:"flavor-id"`
	AvailabilityZone      string              `gcfg:"availability-zone"`
	EnableIngressHostname bool                `gcfg:"enable-ingress-hostname"` // Used with proxy protocol by adding a dns suffix to the load balancer IP address. Default false.
	IngressHostnameSuffix string              `gcfg:"ingress-hostname-suffix"` // Used with proxy protocol by adding a dns suffix to the load balancer IP address. Default nip.io.
	MaxSharedLB           int                 `gcfg:"max-shared-lb"`           //  Number of Services in maximum can share a single load balancer. Default 2
	ContainerStore        string              `gcfg:"container-store"`         // Used to specify the store of the tls-container-ref
	// revive:disable:var-naming
	TlsContainerRef string `gcfg:"default-tls-container-ref"` //  reference to a tls container
	// revive:enable:var-naming
}

// LBClass defines the corresponding floating network, floating subnet or internal subnet ID
type LBClass struct {
	FloatingNetworkID  string `gcfg:"floating-network-id,omitempty"`
	FloatingSubnetID   string `gcfg:"floating-subnet-id,omitempty"`
	FloatingSubnet     string `gcfg:"floating-subnet,omitempty"`
	FloatingSubnetTags string `gcfg:"floating-subnet-tags,omitempty"`
	NetworkID          string `gcfg:"network-id,omitempty"`
	SubnetID           string `gcfg:"subnet-id,omitempty"`
	MemberSubnetID     string `gcfg:"member-subnet-id,omitempty"`
}

func (lb *LBClass) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {
	klog.Info("GetLoadBalancer")
	return &v1.LoadBalancerStatus{}, false, nil
}

// GetLoadBalancerName returns the name of the load balancer. Implementations must treat the
// *v1.Service parameter as read-only and not modify it.
func (lb *LBClass) GetLoadBalancerName(ctx context.Context, clusterName string, service *v1.Service) string {
	klog.Info("GetLoadBalancerName")
	return "ucloudstack"
}

// EnsureLoadBalancer creates a new load balancer 'name', or updates the existing one. Returns the status of the balancer
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (lb *LBClass) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	klog.Info("EnsureLoadBalancer")
	return &v1.LoadBalancerStatus{}, nil
}

// UpdateLoadBalancer updates hosts under the specified load balancer.
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (lb *LBClass) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	klog.Info("UpdateLoadBalancer")
	return nil
}

// EnsureLoadBalancerDeleted deletes the specified load balancer if it
// exists, returning nil if the load balancer specified either didn't exist or
// was successfully deleted.
// This construction is useful because many cloud providers' load balancers
// have multiple underlying components, meaning a Get could say that the LB
// doesn't exist even if some part of it is still laying around.
// Implementations must treat the *v1.Service parameter as read-only and not modify it.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (lb *LBClass) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	klog.Info("EnsureLoadBalancerDeleted")
	return nil
}
