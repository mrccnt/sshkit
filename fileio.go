// Copyright 2017 Marco Conti
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sshkit

// https://stackoverflow.com/questions/28969455/golang-properly-instantiate-os-filemode

// File permission attributes (linux)
const (
	OsRead       = 04
	OsWrite      = 02
	OsEx         = 01
	OsUserShift  = 6
	OsGroupShift = 3
	OsOthShift   = 0

	OsUserR   = OsRead << OsUserShift
	OsUserW   = OsWrite << OsUserShift
	OsUserX   = OsEx << OsUserShift
	OsUserRw  = OsUserR | OsUserW
	OsUserRwx = OsUserRw | OsUserX

	OsGroupR   = OsRead << OsGroupShift
	OsGroupW   = OsWrite << OsGroupShift
	OsGroupX   = OsEx << OsGroupShift
	OsGroupRw  = OsGroupR | OsGroupW
	OsGroupRwx = OsGroupRw | OsGroupX

	OsOthR   = OsRead << OsOthShift
	OsOthW   = OsWrite << OsOthShift
	OsOthX   = OsEx << OsOthShift
	OsOthRw  = OsOthR | OsOthW
	OsOthRwx = OsOthRw | OsOthX

	OsAllR   = OsUserR | OsGroupR | OsOthR
	OsAllW   = OsUserW | OsGroupW | OsOthW
	OsAllX   = OsUserX | OsGroupX | OsOthX
	OsAllRw  = OsAllR | OsAllW
	OsAllRwx = OsAllRw | OsGroupX
)
