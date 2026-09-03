import { haversineDistance, isWithinThreshold } from '../distance';

describe('distance utils', () => {
  it('should calculate distance between two points', () => {
    const lat1 = 19.4326;
    const lng1 = -99.1332;
    const lat2 = 19.4327;
    const lng2 = -99.1333;

    const distance = haversineDistance(lat1, lng1, lat2, lng2);

    expect(typeof distance).toBe('number');
    expect(distance).toBeGreaterThan(0);
    expect(distance).toBeLessThan(100);
  });

  it('should return 0 for same point', () => {
    const lat = 19.4326;
    const lng = -99.1332;

    const distance = haversineDistance(lat, lng, lat, lng);

    expect(distance).toBe(0);
  });

  it('should handle large distances', () => {
    const mexicoCity = { lat: 19.4326, lng: -99.1332 };
    const newYork = { lat: 40.7128, lng: -74.006 };

    const distance = haversineDistance(
      mexicoCity.lat,
      mexicoCity.lng,
      newYork.lat,
      newYork.lng
    );

    expect(distance).toBeGreaterThan(2000000);
    expect(distance).toBeLessThan(4000000);
  });

  it('isWithinThreshold should return true when within threshold', () => {
    expect(isWithinThreshold(100, 500)).toBe(true);
    expect(isWithinThreshold(500, 500)).toBe(true);
  });

  it('isWithinThreshold should return false when outside threshold', () => {
    expect(isWithinThreshold(600, 500)).toBe(false);
  });

  it('isWithinThreshold should use default 500m threshold', () => {
    expect(isWithinThreshold(499)).toBe(true);
    expect(isWithinThreshold(501)).toBe(false);
  });
});
