package privilege

// Kind defines a privilege. This is output by the parser, and used to
// generate the privilege bitfields in the PrivilegeDescriptor.
type Kind uint32

// List of privileges. ALL is specifically encoded so that it will automatically
// pick up new privileges.
// Do not change values of privileges. These correspond to the position
// of the privilege in a bit field and are expected to stay constant.
const (
	ALL    Kind = 1
	CREATE Kind = 2
	DROP   Kind = 3
	// This is a placeholder to make sure that 4 is not reused.
	//
	// It was previously used for the GRANT privilege that has been replaced
	// with the more granular Privilege.GrantOption.
	_                        Kind = 4
	SELECT                   Kind = 5
	INSERT                   Kind = 6
	DELETE                   Kind = 7
	UPDATE                   Kind = 8
	USAGE                    Kind = 9
	ZONECONFIG               Kind = 10
	CONNECT                  Kind = 11
	RULE                     Kind = 12
	MODIFYCLUSTERSETTING     Kind = 13
	EXTERNALCONNECTION       Kind = 14
	VIEWACTIVITY             Kind = 15
	VIEWACTIVITYREDACTED     Kind = 16
	VIEWCLUSTERSETTING       Kind = 17
	CANCELQUERY              Kind = 18
	NOSQLLOGIN               Kind = 19
	EXECUTE                  Kind = 20
	VIEWCLUSTERMETADATA      Kind = 21
	VIEWDEBUG                Kind = 22
	BACKUP                   Kind = 23
	RESTORE                  Kind = 24
	EXTERNALIOIMPLICITACCESS Kind = 25
	CHANGEFEED               Kind = 26
	VIEWJOB                  Kind = 27
	MODIFYSQLCLUSTERSETTING  Kind = 28
	REPLICATION              Kind = 29
	MANAGEVIRTUALCLUSTER     Kind = 30
	VIEWSYSTEMTABLE          Kind = 31
	CREATEROLE               Kind = 32
	CREATELOGIN              Kind = 33
	CREATEDB                 Kind = 34
	CONTROLJOB               Kind = 35
	REPAIRCLUSTERMETADATA    Kind = 36
	largestKind                   = REPAIRCLUSTERMETADATA
)
