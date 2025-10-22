import { useEffect, useState } from "react";
import { MapContainer, TileLayer, Marker, Popup } from "react-leaflet";
import { Button } from "@/components/ui/button";
import { MapPin, Utensils, Hotel, Plane } from "lucide-react";
import L from "leaflet";

// Fix for default markers in react-leaflet
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png",
  iconUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png",
  shadowUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png",
});

interface Location {
  name: string;
  lat: number;
  lng: number;
  type: string;
}

interface MapPanelProps {
  locations: Location[];
}

export function MapPanel({ locations }: MapPanelProps) {
  const [center, setCenter] = useState<[number, number]>([12.3375, 75.8069]); // Default to Coorg

  useEffect(() => {
    if (locations.length > 0) {
      const mainLocation =
        locations.find((loc) => loc.type === "destination") || locations[0];
      setCenter([mainLocation.lat, mainLocation.lng]);
    }
  }, [locations]);

  return (
    <div className="h-full flex flex-col bg-white">
      {/* Itinerary Header */}
      <div className="mb-4">
        <h2 className="text-xl font-bold text-black mb-2">
          3-Day Itinerary in Coorg
        </h2>
        <div className="flex items-center justify-between text-sm">
          <span className="text-gray-500">Feb 12 - Feb 19, 2025</span>
          <span className="text-purple-600 font-semibold">
            Total Cost: 10000Rs
          </span>
        </div>
      </div>

      {/* Map */}
      <div className="mb-4 comic-border rounded-2xl overflow-hidden">
        <div className="h-48 w-full">
          <MapContainer
            center={center}
            zoom={11}
            style={{ height: "100%", width: "100%" }}
            className="rounded-2xl"
          >
            <TileLayer
              attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />
            {locations.map((location, index) => (
              <Marker key={index} position={[location.lat, location.lng]}>
                <Popup>
                  <div>
                    <strong>{location.name}</strong>
                    <br />
                    Type: {location.type}
                  </div>
                </Popup>
              </Marker>
            ))}
          </MapContainer>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="flex space-x-2 mb-4">
        <Button
          variant="outline"
          size="sm"
          className="flex-1 comic-border rounded-xl text-xs"
        >
          <MapPin className="w-3 h-3 mr-1" />
          Itinerary
        </Button>
        <Button
          variant="outline"
          size="sm"
          className="flex-1 comic-border rounded-xl text-xs"
        >
          <Utensils className="w-3 h-3 mr-1" />
          Restaurant
        </Button>
        <Button
          variant="outline"
          size="sm"
          className="flex-1 comic-border rounded-xl text-xs"
        >
          <Hotel className="w-3 h-3 mr-1" />
          Hotel
        </Button>
        <Button
          variant="outline"
          size="sm"
          className="flex-1 comic-border rounded-xl text-xs"
        >
          <Plane className="w-3 h-3 mr-1" />
          Flights
        </Button>
      </div>

      {/* Day Tabs */}
      <div className="flex space-x-2 mb-4">
        <Button className="flex-1 comic-button rounded-xl text-sm">
          Day 1
        </Button>
        <Button
          variant="outline"
          className="flex-1 comic-border rounded-xl text-sm"
        >
          Day 2
        </Button>
        <Button
          variant="outline"
          className="flex-1 comic-border rounded-xl text-sm"
        >
          Day 3
        </Button>
      </div>

      {/* Itinerary Details */}
      <div className="flex-1 overflow-y-auto">
        <div className="space-y-3">
          <div className="flex items-start space-x-3 p-3 rounded-xl hover:bg-gray-50">
            <div className="w-8 h-8 bg-purple-100 rounded-full flex items-center justify-center flex-shrink-0">
              <span className="text-sm font-semibold text-purple-600">1</span>
            </div>
            <div className="flex-1">
              <h4 className="font-semibold text-sm text-black">Pera Museum</h4>
              <p className="text-xs text-gray-600">4.7 (441 reviews)</p>
              <p className="text-xs text-gray-500 mt-1">Wednesday, Feb 10</p>
            </div>
            <div className="w-12 h-12 bg-gray-200 rounded-lg flex-shrink-0"></div>
          </div>

          <div className="flex items-start space-x-3 p-3 rounded-xl hover:bg-gray-50">
            <div className="w-8 h-8 bg-purple-100 rounded-full flex items-center justify-center flex-shrink-0">
              <span className="text-sm font-semibold text-purple-600">2</span>
            </div>
            <div className="flex-1">
              <h4 className="font-semibold text-sm text-black">Pera Museum</h4>
              <p className="text-xs text-gray-600">4.7 (441 reviews)</p>
              <p className="text-xs text-gray-500 mt-1">Wednesday, Feb 10</p>
            </div>
            <div className="w-12 h-12 bg-gray-200 rounded-lg flex-shrink-0"></div>
          </div>

          <div className="flex items-start space-x-3 p-3 rounded-xl hover:bg-gray-50">
            <div className="w-8 h-8 bg-purple-100 rounded-full flex items-center justify-center flex-shrink-0">
              <span className="text-sm font-semibold text-purple-600">3</span>
            </div>
            <div className="flex-1">
              <h4 className="font-semibold text-sm text-black">Pera Museum</h4>
              <p className="text-xs text-gray-600">4.7 (441 reviews)</p>
              <p className="text-xs text-gray-500 mt-1">Wednesday, Feb 10</p>
            </div>
            <div className="w-12 h-12 bg-gray-200 rounded-lg flex-shrink-0"></div>
          </div>
        </div>
      </div>
    </div>
  );
}
