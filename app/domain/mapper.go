// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package domain

// Mapper maps fields from one domain model to another. These models can be
// understood as data transfer objects (DTOs). DOTs here are used to keep
// data sanitized between adapters and the service layer. For instance, it
// may ensure that sensitive data returned from the service layer, such as a
// password hash does not make it to transmission at the port API layer.
//
// Mappers are often used at the boundary between port and adapter, and not at
// the service layer.
type Mapper struct{}
