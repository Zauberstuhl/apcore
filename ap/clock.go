// apcore is a server framework for implementing an ActivityPub application.
// Copyright (C) 2019 Cory Slep
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package ap

import (
	"time"

	"github.com/go-fed/activity/pub"
)

var _ pub.Clock = &Clock{}

type Clock struct {
	loc *time.Location
}

// Creates new clock with IANA Time Zone database string
func NewClock(location string) (c *Clock, err error) {
	c = &Clock{}
	c.loc, err = time.LoadLocation(location)
	return
}

func (c *Clock) Now() time.Time {
	return time.Now().In(c.loc)
}
