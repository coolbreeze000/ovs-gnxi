
/* Copyright (c) 2015 Open Networking Foundation.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at:
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#ifndef OFC_COMM_SOCK_H_
#define OFC_COMM_SOCK_H_

#include <config.h>

#define OFC_SOCK_SENDFLAGS MSG_NOSIGNAL

#define OFC_SOCK_PATH OFC_DATADIR"/ofc.sock"
#ifdef OFC_SOCK_GROUP
#	define OFC_SOCK_PERM 0660
#else
#	define OFC_SOCK_PERM 0666
#endif

/* communication handler */
typedef int comm_t;




typedef int msgtype_t;

enum COMM_SOCKET_MSGTYPE {
	COMM_SOCK_RESULT_ERROR = -1,
	COMM_SOCK_GET_CPBLTS = 1,
	COMM_SOCK_SET_SESSION,
	COMM_SOCK_CLOSE_SESSION,
	COMM_SOCK_KILL_SESSION,
	COMM_SOCK_GENERICOP
};

#endif /* OFC_COMM_SOCK_H_ */
