// Copyright (C) 2019-2022 Algorand, Inc.
// This file is part of go-algorand
//
// go-algorand is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// go-algorand is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with go-algorand.  If not, see <https://www.gnu.org/licenses/>.

package config

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/algorand/go-algorand/protocol"
	"github.com/algorand/go-algorand/util/codecs"
)

// Local holds the per-node-instance configuration settings for the protocol.
// !!! WARNING !!!
//
// These versioned struct tags need to be maintained CAREFULLY and treated
// like UNIVERSAL CONSTANTS - they should not be modified once committed.
//
// New fields may be added to the Local struct, along with a version tag
// denoting a new version. When doing so, also update the
// test/testdata/configs/config-v{n}.json and call "make generate" to regenerate the constants.
//
// !!! WARNING !!!
type Local struct {
	// Version tracks the current version of the defaults so we can migrate old -> new
	// This is specifically important whenever we decide to change the default value
	// for an existing parameter. This field tag must be updated any time we add a new version.
	Version uint32 `version[0]:"0" version[1]:"1" version[2]:"2" version[3]:"3" version[4]:"4" version[5]:"5" version[6]:"6" version[7]:"7" version[8]:"8" version[9]:"9" version[10]:"10" version[11]:"11" version[12]:"12" version[13]:"13" version[14]:"14" version[15]:"15" version[16]:"16" version[17]:"17" version[18]:"18" version[19]:"19" version[20]:"20" version[21]:"21" version[22]:"22"`

	// environmental (may be overridden)
	// When enabled, stores blocks indefinitally, otherwise, only the most recents blocks
	// are being kept around. ( the precise number of recent blocks depends on the consensus parameters )
	Archival bool `version[0]:"false"`

	// gossipNode.go
	// how many peers to propagate to?
	GossipFanout int    `version[0]:"4"`
	NetAddress   string `version[0]:""`

	// 1 * time.Minute = 60000000000 ns
	ReconnectTime time.Duration `version[0]:"60" version[1]:"60000000000"`

	// what we should tell people to connect to
	PublicAddress string `version[0]:""`

	MaxConnectionsPerIP int `version[3]:"30"`

	// 0 == disable
	PeerPingPeriodSeconds int `version[0]:"0"`

	// for https serving
	TLSCertFile string `version[0]:""`
	TLSKeyFile  string `version[0]:""`

	// Logging
	BaseLoggerDebugLevel uint32 `version[0]:"1" version[1]:"4"`
	// if this is 0, do not produce agreement.cadaver
	CadaverSizeTarget uint64 `version[0]:"1073741824"`

	// IncomingConnectionsLimit specifies the max number of long-lived incoming
	// connections. 0 means no connections allowed. Must be non-negative.
	// Estimating 5MB per incoming connection, 5MB*800 = 4GB
	IncomingConnectionsLimit int `version[0]:"-1" version[1]:"10000" version[17]:"800"`

	// BroadcastConnectionsLimit specifies the number of connections that
	// will receive broadcast (gossip) messages from this node.  If the
	// node has more connections than this number, it will send broadcasts
	// to the top connections by priority (outgoing connections first, then
	// by money held by peers based on their participation key).  0 means
	// no outgoing messages (not even transaction broadcasting to outgoing
	// peers).  -1 means unbounded (default).
	BroadcastConnectionsLimit int `version[4]:"-1"`

	// AnnounceParticipationKey specifies that this node should announce its
	// participation key (with the largest stake) to its gossip peers.  This
	// allows peers to prioritize our connection, if necessary, in case of a
	// DoS attack.  Disabling this means that the peers will not have any
	// additional information to allow them to prioritize our connection.
	AnnounceParticipationKey bool `version[4]:"true"`

	// PriorityPeers specifies peer IP addresses that should always get
	// outgoing broadcast messages from this node.
	PriorityPeers map[string]bool `version[4]:""`

	// To make sure the algod process does not run out of FDs, algod ensures
	// that RLIMIT_NOFILE >= IncomingConnectionsLimit + RestConnectionsHardLimit +
	// ReservedFDs. ReservedFDs are meant to leave room for short-lived FDs like
	// DNS queries, SQLite files, etc. This parameter shouldn't be changed.
	ReservedFDs uint64 `version[2]:"256"`

	// local server
	// API endpoint address
	EndpointAddress string `version[0]:"127.0.0.1:0"`

	// timeouts passed to the rest http.Server implementation
	RestReadTimeoutSeconds  int `version[4]:"15"`
	RestWriteTimeoutSeconds int `version[4]:"120"`

	// SRV-based phonebook
	DNSBootstrapID string `version[0]:"<network>.algorand.network"`

	// Log file size limit in bytes. When set to 0 logs will be written to stdout.
	LogSizeLimit uint64 `version[0]:"1073741824"`

	// text/template for creating log archive filename.
	// Available template vars:
	// Time at start of log: {{.Year}} {{.Month}} {{.Day}} {{.Hour}} {{.Minute}} {{.Second}}
	// Time at end of log: {{.EndYear}} {{.EndMonth}} {{.EndDay}} {{.EndHour}} {{.EndMinute}} {{.EndSecond}}
	//
	// If the filename ends with .gz or .bz2 it will be compressed.
	//
	// default: "node.archive.log" (no rotation, clobbers previous archive)
	LogArchiveName string `version[4]:"node.archive.log"`

	// LogArchiveMaxAge will be parsed by time.ParseDuration().
	// Valid units are 's' seconds, 'm' minutes, 'h' hours
	LogArchiveMaxAge string `version[4]:""`

	// number of consecutive attempts to catchup after which we replace the peers we're connected to
	CatchupFailurePeerRefreshRate int `version[0]:"10"`

	// where should the node exporter listen for metrics
	NodeExporterListenAddress string `version[0]:":9100"`

	// enable metric reporting flag
	EnableMetricReporting bool `version[0]:"false"`

	// enable top accounts reporting flag
	EnableTopAccountsReporting bool `version[0]:"false"`

	// enable agreement reporting flag. Currently only prints additional period events.
	EnableAgreementReporting bool `version[3]:"false"`

	// enable agreement timing metrics flag
	EnableAgreementTimeMetrics bool `version[3]:"false"`

	// The path to the node exporter.
	NodeExporterPath string `version[0]:"./node_exporter"`

	// The fallback DNS resolver address that would be used if the system resolver would fail to retrieve SRV records
	FallbackDNSResolverAddress string `version[0]:""`

	// exponential increase factor of transaction pool's fee threshold, should always be 2 in production
	TxPoolExponentialIncreaseFactor uint64 `version[0]:"2"`

	SuggestedFeeBlockHistory int `version[0]:"3"`

	// TxPoolSize is the number of transactions that fit in the transaction pool
	TxPoolSize int `version[0]:"50000" version[5]:"15000"`

	// number of seconds allowed for syncing transactions
	TxSyncTimeoutSeconds int64 `version[0]:"30"`

	// number of seconds between transaction synchronizations
	TxSyncIntervalSeconds int64 `version[0]:"60"`

	// the number of incoming message hashes buckets.
	IncomingMessageFilterBucketCount int `version[0]:"5"`

	// the size of each incoming message hash bucket.
	IncomingMessageFilterBucketSize int `version[0]:"512"`

	// the number of outgoing message hashes buckets.
	OutgoingMessageFilterBucketCount int `version[0]:"3"`

	// the size of each outgoing message hash bucket.
	OutgoingMessageFilterBucketSize int `version[0]:"128"`

	// enable the filtering of outgoing messages
	EnableOutgoingNetworkMessageFiltering bool `version[0]:"true"`

	// enable the filtering of incoming messages
	EnableIncomingMessageFilter bool `version[0]:"false"`

	// control enabling / disabling deadlock detection.
	// negative (-1) to disable, positive (1) to enable, 0 for default.
	DeadlockDetection int `version[1]:"0"`

	// The threshold used for deadlock detection, in seconds.
	DeadlockDetectionThreshold int `version[20]:"30"`

	// Prefer to run algod Hosted (under algoh)
	// Observed by `goal` for now.
	RunHosted bool `version[3]:"false"`

	// The maximal number of blocks that catchup will fetch in parallel.
	// If less than Protocol.SeedLookback, then Protocol.SeedLookback will be used as to limit the catchup.
	// Setting this variable to 0 would disable the catchup
	CatchupParallelBlocks uint64 `version[3]:"50" version[5]:"16"`

	// Generate AssembleBlockMetrics telemetry event
	EnableAssembleStats bool `version[0]:""`

	// Generate ProcessBlockMetrics telemetry event
	EnableProcessBlockStats bool `version[0]:""`

	// SuggestedFeeSlidingWindowSize is number of past blocks that will be considered in computing the suggested fee
	SuggestedFeeSlidingWindowSize uint32 `version[3]:"50"`

	// the max size the sync server would return
	TxSyncServeResponseSize int `version[3]:"1000000"`

	// IsIndexerActive indicates whether to activate the indexer for fast retrieval of transactions
	// Note -- Indexer cannot operate on non Archival nodes
	IsIndexerActive bool `version[3]:"false"`

	// UseXForwardedForAddress indicates whether or not the node should use the X-Forwarded-For HTTP Header when
	// determining the source of a connection.  If used, it should be set to the string "X-Forwarded-For", unless the
	// proxy vendor provides another header field.  In the case of CloudFlare proxy, the "CF-Connecting-IP" header
	// field can be used.
	UseXForwardedForAddressField string `version[0]:""`

	// ForceRelayMessages indicates whether the network library relay messages even in the case that no NetAddress was specified.
	ForceRelayMessages bool `version[0]:"false"`

	// ConnectionsRateLimitingWindowSeconds is being used in conjunction with ConnectionsRateLimitingCount;
	// see ConnectionsRateLimitingCount description for further information. Providing a zero value
	// in this variable disables the connection rate limiting.
	ConnectionsRateLimitingWindowSeconds uint `version[4]:"1"`

	// ConnectionsRateLimitingCount is being used along with ConnectionsRateLimitingWindowSeconds to determine if
	// a connection request should be accepted or not. The gossip network examine all the incoming requests in the past
	// ConnectionsRateLimitingWindowSeconds seconds that share the same origin. If the total count exceed the ConnectionsRateLimitingCount
	// value, the connection is refused.
	ConnectionsRateLimitingCount uint `version[4]:"60"`

	// EnableRequestLogger enabled the logging of the incoming requests to the telemetry server.
	EnableRequestLogger bool `version[4]:"false"`

	// PeerConnectionsUpdateInterval defines the interval at which the peer connections information is being sent to the
	// telemetry ( when enabled ). Defined in seconds.
	PeerConnectionsUpdateInterval int `version[5]:"3600"`

	// EnableProfiler enables the go pprof endpoints, should be false if
	// the algod api will be exposed to untrusted individuals
	EnableProfiler bool `version[0]:"false"`

	// EnableRuntimeMetrics exposes Go runtime metrics in /metrics and via node_exporter.
	EnableRuntimeMetrics bool `version[22]:"false"`

	// TelemetryToLog records messages to node.log that are normally sent to remote event monitoring
	TelemetryToLog bool `version[5]:"true"`

	// DNSSecurityFlags instructs algod validating DNS responses.
	// Possible fla values
	// 0x00 - disabled
	// 0x01 (dnssecSRV) - validate SRV response
	// 0x02 (dnssecRelayAddr) - validate relays' names to addresses resolution
	// 0x04 (dnssecTelemetryAddr) - validate telemetry and metrics names to addresses resolution
	// ...
	DNSSecurityFlags uint32 `version[6]:"1"`

	// EnablePingHandler controls whether the gossip node would respond to ping messages with a pong message.
	EnablePingHandler bool `version[6]:"true"`

	// DisableOutgoingConnectionThrottling disables the connection throttling of the network library, which
	// allow the network library to continuesly disconnect relays based on their relative ( and absolute ) performance.
	DisableOutgoingConnectionThrottling bool `version[5]:"false"`

	// NetworkProtocolVersion overrides network protocol version ( if present )
	NetworkProtocolVersion string `version[6]:""`

	// CatchpointInterval sets the interval at which catchpoint are being generated. Setting this to 0 disables the catchpoint from being generated.
	// See CatchpointTracking for more details.
	CatchpointInterval uint64 `version[7]:"10000"`

	// CatchpointFileHistoryLength defines how many catchpoint files we want to store back.
	// 0 means don't store any, -1 mean unlimited and positive number suggest the number of most recent catchpoint files.
	CatchpointFileHistoryLength int `version[7]:"365"`

	// EnableLedgerService enables the ledger serving service. The functionality of this depends on NetAddress, which must also be provided.
	// This functionality is required for the catchpoint catchup.
	EnableLedgerService bool `version[7]:"false"`

	// EnableBlockService enables the block serving service. The functionality of this depends on NetAddress, which must also be provided.
	// This functionality is required for the catchup.
	EnableBlockService bool `version[7]:"false"`

	// EnableGossipBlockService enables the block serving service over the gossip network. The functionality of this depends on NetAddress, which must also be provided.
	// This functionality is required for the relays to perform catchup from nodes.
	EnableGossipBlockService bool `version[8]:"true"`

	// CatchupHTTPBlockFetchTimeoutSec controls how long the http query for fetching a block from a relay would take before giving up and trying another relay.
	CatchupHTTPBlockFetchTimeoutSec int `version[9]:"4"`

	// CatchupGossipBlockFetchTimeoutSec controls how long the gossip query for fetching a block from a relay would take before giving up and trying another relay.
	CatchupGossipBlockFetchTimeoutSec int `version[9]:"4"`

	// CatchupLedgerDownloadRetryAttempts controls the number of attempt the ledger fetching would be attempted before giving up catching up to the provided catchpoint.
	CatchupLedgerDownloadRetryAttempts int `version[9]:"50"`

	// CatchupLedgerDownloadRetryAttempts controls the number of attempt the block fetching would be attempted before giving up catching up to the provided catchpoint.
	CatchupBlockDownloadRetryAttempts int `version[9]:"1000"`

	// EnableDeveloperAPI enables teal/compile, teal/dryrun API endpoints.
	// This functionality is disabled by default.
	EnableDeveloperAPI bool `version[9]:"false"`

	// OptimizeAccountsDatabaseOnStartup controls whether the accounts database would be optimized
	// on algod startup.
	OptimizeAccountsDatabaseOnStartup bool `version[10]:"false"`

	// CatchpointTracking determines if catchpoints are going to be tracked. The value is interpreted as follows:
	// A value of -1 means "don't track catchpoints".
	// A value of 1 means "track catchpoints as long as CatchpointInterval is also set to a positive non-zero value". If CatchpointInterval <= 0, no catchpoint tracking would be performed.
	// A value of 0 means automatic, which is the default value. In this mode, a non archival node would not track the catchpoints, and an archival node would track the catchpoints as long as CatchpointInterval > 0.
	// Other values of CatchpointTracking would give a warning in the log file, and would behave as if the default value was provided.
	CatchpointTracking int64 `version[11]:"0"`

	// LedgerSynchronousMode defines the synchronous mode used by the ledger database. The supported options are:
	// 0 - SQLite continues without syncing as soon as it has handed data off to the operating system.
	// 1 - SQLite database engine will still sync at the most critical moments, but less often than in FULL mode.
	// 2 - SQLite database engine will use the xSync method of the VFS to ensure that all content is safely written to the disk surface prior to continuing. On Mac OS, the data is additionally syncronized via fullfsync.
	// 3 - In addition to what being done in 2, it provides additional durability if the commit is followed closely by a power loss.
	// for further information see the description of SynchronousMode in dbutil.go
	LedgerSynchronousMode int `version[12]:"2"`

	// AccountsRebuildSynchronousMode defines the synchronous mode used by the ledger database while the account database is being rebuilt. This is not a typical operational usecase,
	// and is expected to happen only on either startup ( after enabling the catchpoint interval, or on certain database upgrades ) or during fast catchup. The values specified here
	// and their meanings are identical to the ones in LedgerSynchronousMode.
	AccountsRebuildSynchronousMode int `version[12]:"1"`

	// MaxCatchpointDownloadDuration defines the maximum duration a client will be keeping the outgoing connection of a catchpoint download request open for processing before
	// shutting it down. Networks that have large catchpoint files, slow connection or slow storage could be a good reason to increase this value. Note that this is a client-side only
	// configuration value, and it's independent of the actual catchpoint file size.
	MaxCatchpointDownloadDuration time.Duration `version[13]:"7200000000000"`

	// MinCatchpointFileDownloadBytesPerSecond defines the minimal download speed that would be considered to be "acceptable" by the catchpoint file fetcher, measured in bytes per seconds. If the
	// provided stream speed drops below this threshold, the connection would be recycled. Note that this field is evaluated per catchpoint "chunk" and not on it's own. If this field is zero,
	// the default of 20480 would be used.
	MinCatchpointFileDownloadBytesPerSecond uint64 `version[13]:"20480"`

	// TraceServer is a host:port to report graph propagation trace info to.
	NetworkMessageTraceServer string `version[13]:""`

	// VerifiedTranscationsCacheSize defines the number of transactions that the verified transactions cache would hold before cycling the cache storage in a round-robin fashion.
	VerifiedTranscationsCacheSize int `version[14]:"30000"`

	// EnableCatchupFromArchiveServers controls which peers the catchup service would use in order to catchup.
	// When enabled, the catchup service would use the archive servers before falling back to the relays.
	// On networks that doesn't have archive servers, this becomes a no-op, as the catchup service would have no
	// archive server to pick from, and therefore automatically selects one of the relay nodes.
	EnableCatchupFromArchiveServers bool `version[15]:"false"`

	// DisableLocalhostConnectionRateLimit controls whether the incoming connection rate limit would apply for
	// connections that are originating from the local machine. Setting this to "true", allow to create large
	// local-machine networks that won't trip the incoming connection limit observed by relays.
	DisableLocalhostConnectionRateLimit bool `version[16]:"true"`

	// BlockServiceCustomFallbackEndpoints is a comma delimited list of endpoints which the block service uses to
	// redirect the http requests to in case it does not have the round. If it is not specified, will check
	// EnableBlockServiceFallbackToArchiver.
	BlockServiceCustomFallbackEndpoints string `version[16]:""`

	// EnableBlockServiceFallbackToArchiver controls whether the block service redirects the http requests to
	// an archiver or return StatusNotFound (404) when in does not have the requested round, and
	// BlockServiceCustomFallbackEndpoints is empty.
	// The archiver is randomly selected, if none is available, will return StatusNotFound (404).
	EnableBlockServiceFallbackToArchiver bool `version[16]:"true"`

	// CatchupBlockValidateMode is a development and testing configuration used by the catchup service.
	// It can be used to omit certain validations to speed up the catchup process, or to apply extra validations which are redundant in normal operation.
	// This field is a bit-field with:
	// bit 0: (default 0) 0: verify the block certificate; 1: skip this validation
	// bit 1: (default 0) 0: verify payset committed hash in block header matches payset hash; 1: skip this validation
	// bit 2: (default 0) 0: don't verify the transaction signatures on the block are valid; 1: verify the transaction signatures on block
	// bit 3: (default 0) 0: don't verify that the hash of the recomputed payset matches the hash of the payset committed in the block header; 1: do perform the above verification
	// Note: not all permutations of the above bitset are currently functional. In particular, the ones that are functional are:
	// 0  : default behavior.
	// 3  : speed up catchup by skipping necessary validations
	// 12 : perform all validation methods (normal and additional). These extra tests helps to verify the integrity of the compiled executable against
	//      previously used executabled, and would not provide any additional security guarantees.
	CatchupBlockValidateMode int `version[16]:"0"`

	// Generate AccountUpdates telemetry event
	EnableAccountUpdatesStats bool `version[16]:"false"`

	// Time interval in nanoseconds for generating accountUpdates telemetry event
	AccountUpdatesStatsInterval time.Duration `version[16]:"5000000000"`

	// ParticipationKeysRefreshInterval is the duration between two consecutive checks to see if new participation
	// keys have been placed on the genesis directory.
	ParticipationKeysRefreshInterval time.Duration `version[16]:"60000000000"`

	// DisableNetworking disables all the incoming and outgoing communication a node would perform. This is useful
	// when we have a single-node private network, where there is no other nodes that need to be communicated with.
	// features like catchpoint catchup would be rendered completly non-operational, and many of the node inner
	// working would be completly dis-functional.
	DisableNetworking bool `version[16]:"false"`

	// ForceFetchTransactions allows to explicitly configure a node to retrieve all the transactions
	// into it's transaction pool, even if those would not be required as the node doesn't
	// participate in the consensus or used to relay transactions.
	ForceFetchTransactions bool `version[17]:"false"`

	// EnableVerbosedTransactionSyncLogging enables the transaction sync to write extensive
	// message exchange information to the log file. This option is disabled by default,
	// so that the log files would not grow too rapidly.
	EnableVerbosedTransactionSyncLogging bool `version[17]:"false"`

	// TransactionSyncDataExchangeRate overrides the auto-calculated data exchange rate between each
	// two peers. The unit of the data exchange rate is in bytes per second. Setting the value to
	// zero implies allowing the transaction sync to dynamically calculate the value.
	TransactionSyncDataExchangeRate uint64 `version[17]:"0"`

	// TransactionSyncSignificantMessageThreshold define the threshold used for a transaction sync
	// message before it can be used for calculating the data exchange rate. Setting this to zero
	// would use the default values. The threshold is defined in units of bytes.
	TransactionSyncSignificantMessageThreshold uint64 `version[17]:"0"`

	// ProposalAssemblyTime is the max amount of time to spend on generating a proposal block.
	ProposalAssemblyTime time.Duration `version[19]:"250000000"`

	// When the number of http connections to the REST layer exceeds the soft limit,
	// we start returning http code 429 Too Many Requests.
	RestConnectionsSoftLimit uint64 `version[20]:"1024"`
	// The http server does not accept new connections as long we have this many
	// (hard limit) connections already.
	RestConnectionsHardLimit uint64 `version[20]:"2048"`

	// MaxAPIResourcesPerAccount sets the maximum total number of resources (created assets, created apps,
	// asset holdings, and application local state) per account that will be allowed in AccountInformation
	// REST API responses before returning a 400 Bad Request. Set zero for no limit.
	MaxAPIResourcesPerAccount uint64 `version[21]:"100000"`

	// AgreementIncomingVotesQueueLength sets the size of the buffer holding incoming votes.
	AgreementIncomingVotesQueueLength uint64 `version[21]:"10000"`

	// AgreementIncomingProposalsQueueLength sets the size of the buffer holding incoming proposals.
	AgreementIncomingProposalsQueueLength uint64 `version[21]:"25"`

	// AgreementIncomingBundlesQueueLength sets the size of the buffer holding incoming bundles.
	AgreementIncomingBundlesQueueLength uint64 `version[21]:"7"`
}

// DNSBootstrapArray returns an array of one or more DNS Bootstrap identifiers
func (cfg Local) DNSBootstrapArray(networkID protocol.NetworkID) (bootstrapArray []string) {
	dnsBootstrapString := cfg.DNSBootstrap(networkID)
	bootstrapArray = strings.Split(dnsBootstrapString, ";")
	// omit zero length entries from the result set.
	for i := len(bootstrapArray) - 1; i >= 0; i-- {
		if len(bootstrapArray[i]) == 0 {
			bootstrapArray = append(bootstrapArray[:i], bootstrapArray[i+1:]...)
		}
	}
	return
}

// DNSBootstrap returns the network-specific DNSBootstrap identifier
func (cfg Local) DNSBootstrap(network protocol.NetworkID) string {
	// if user hasn't modified the default DNSBootstrapID in the configuration
	// file and we're targeting a devnet ( via genesis file ), we the
	// explicit devnet network bootstrap.
	if defaultLocal.DNSBootstrapID == cfg.DNSBootstrapID {
		if network == Devnet {
			return "devnet.algodev.network"
		} else if network == Betanet {
			return "betanet.algodev.network"
		}
	}
	return strings.Replace(cfg.DNSBootstrapID, "<network>", string(network), -1)
}

// SaveToDisk writes the non-default Local settings into a root/ConfigFilename file
func (cfg Local) SaveToDisk(root string) error {
	configpath := filepath.Join(root, ConfigFilename)
	filename := os.ExpandEnv(configpath)
	return cfg.SaveToFile(filename)
}

// SaveAllToDisk writes the all Local settings into a root/ConfigFilename file
func (cfg Local) SaveAllToDisk(root string) error {
	configpath := filepath.Join(root, ConfigFilename)
	filename := os.ExpandEnv(configpath)
	prettyPrint := true
	return codecs.SaveObjectToFile(filename, cfg, prettyPrint)
}

// SaveToFile saves the config to a specific filename, allowing overriding the default name
func (cfg Local) SaveToFile(filename string) error {
	var alwaysInclude []string
	alwaysInclude = append(alwaysInclude, "Version")
	return codecs.SaveNonDefaultValuesToFile(filename, cfg, defaultLocal, alwaysInclude, true)
}

// DNSSecuritySRVEnforced returns true if SRV response verification enforced
func (cfg Local) DNSSecuritySRVEnforced() bool {
	return cfg.DNSSecurityFlags&dnssecSRV != 0
}

// DNSSecurityRelayAddrEnforced returns true if relay name to ip addr resolution enforced
func (cfg Local) DNSSecurityRelayAddrEnforced() bool {
	return cfg.DNSSecurityFlags&dnssecRelayAddr != 0
}

// DNSSecurityTelemeryAddrEnforced returns true if relay name to ip addr resolution enforced
func (cfg Local) DNSSecurityTelemeryAddrEnforced() bool {
	return cfg.DNSSecurityFlags&dnssecTelemetryAddr != 0
}

// CatchupVerifyCertificate returns true if certificate verification is needed
func (cfg Local) CatchupVerifyCertificate() bool {
	return cfg.CatchupBlockValidateMode&catchupValidationModeCertificate == 0
}

// CatchupVerifyPaysetHash returns true if payset hash verification is needed
func (cfg Local) CatchupVerifyPaysetHash() bool {
	return cfg.CatchupBlockValidateMode&catchupValidationModePaysetHash == 0
}

// CatchupVerifyTransactionSignatures returns true if transactions signature verification is needed
func (cfg Local) CatchupVerifyTransactionSignatures() bool {
	return cfg.CatchupBlockValidateMode&catchupValidationModeVerifyTransactionSignatures != 0
}

// CatchupVerifyApplyData returns true if verifying the ApplyData of the payset needed
func (cfg Local) CatchupVerifyApplyData() bool {
	return cfg.CatchupBlockValidateMode&catchupValidationModeVerifyApplyData != 0
}
