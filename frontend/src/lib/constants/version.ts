/**
 * Kodia Framework Version Information
 * This is automatically updated during releases
 */

export const KODIA_VERSION = '1.0.0'
export const KODIA_RELEASE_TYPE = 'stable'
export const KODIA_BUILD_DATE = '2026-05-21'

/**
 * Version parts for programmatic use
 */
export const VERSION_INFO = {
  major: 1,
  minor: 0,
  patch: 0,
  prerelease: '',
  buildMetadata: '',
}

/**
 * Returns the full version string
 */
export function getKodiaVersion(): string {
  let version = `${VERSION_INFO.major}.${VERSION_INFO.minor}.${VERSION_INFO.patch}`
  if (VERSION_INFO.prerelease) {
    version += `-${VERSION_INFO.prerelease}`
  }
  if (VERSION_INFO.buildMetadata) {
    version += `+${VERSION_INFO.buildMetadata}`
  }
  return version
}

/**
 * Returns version with build information
 */
export function getKodiaVersionDetails(): string {
  const version = getKodiaVersion()
  const status = KODIA_RELEASE_TYPE === 'stable' ? '✅ Stable' : `⚠️ ${KODIA_RELEASE_TYPE}`
  return `Kodia Framework v${version} (${status}, built: ${KODIA_BUILD_DATE})`
}

/**
 * Check if current version is compatible with required version
 * Simple semver compatibility check
 */
export function isVersionCompatible(required: string): boolean {
  const [reqMajor, rest] = required.split('.')
  const reqMinor = rest?.split('.')[0]

  const currentMajor = VERSION_INFO.major.toString()
  const currentMinor = VERSION_INFO.minor.toString()

  // For v1.x.x, require at least v1.0.0
  if (reqMajor === '1') {
    return currentMajor === '1' && parseInt(currentMinor) >= parseInt(reqMinor || '0')
  }

  return currentMajor === reqMajor
}

// Log version info in development
if (typeof window !== 'undefined') {
  if (process.env.NODE_ENV === 'development') {
    console.log(`%c${getKodiaVersionDetails()}`, 'color: #3b82f6; font-weight: bold;')
  }
}
