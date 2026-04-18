"use client";

import React, { useState, useEffect } from "react";
import { Modal } from "@/components/ui/modal";
import Button from "@/components/ui/button/Button";
import Label from "@/components/form/Label";
import Input from "@/components/form/input/InputField";
import TextArea from "@/components/form/input/TextArea";
import Select from "@/components/form/Select";
import type { ServiceOverview, GatewayData, HtpasswdSecretSummary, HTTPRouteCreateRequest, ParentRef } from "@/lib/api";
import { createHTTPRoute, listHtpasswdSecrets } from "@/lib/api";

interface HTTPRouteCreateModalProps {
  isOpen: boolean;
  onClose: () => void;
  service: ServiceOverview;
  gateways: GatewayData[];
  onCreate: () => void;
}

export default function HTTPRouteCreateModal({
  isOpen,
  onClose,
  service,
  gateways,
  onCreate,
}: HTTPRouteCreateModalProps) {
  const [name, setName] = useState("");
  const [hostnames, setHostnames] = useState("");
  const [selectedGateway, setSelectedGateway] = useState("");
  const [port, setPort] = useState<number>(service.ports[0]?.port || 80);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [key, setKey] = useState(0);

  // SecurityPolicy options
  const [createSecurityPolicy, setCreateSecurityPolicy] = useState(false);
  const [htpasswdSecrets, setHtpasswdSecrets] = useState<HtpasswdSecretSummary[]>([]);
  const [selectedSecret, setSelectedSecret] = useState("");

  // Reset form when modal opens
  useEffect(() => {
    if (isOpen) {
      setName("");
      setHostnames("");
      setSelectedGateway("");
      setPort(service.ports[0]?.port || 80);
      setError(null);
      setCreateSecurityPolicy(false);
      setSelectedSecret("");
      setKey((prev) => prev + 1); // Force re-mount of Select

      // Fetch Htpasswd secrets for this namespace
      listHtpasswdSecrets(service.namespace)
        .then((res) => setHtpasswdSecrets(res.data))
        .catch(() => {
          // Htpasswd secrets are optional, ignore error
        });
    }
  }, [isOpen, service]);

  // Prepare gateway options for select
  const gatewayOptions = gateways.map((gw) => ({
    value: `${gw.namespace}/${gw.name}`,
    label: `${gw.namespace}/${gw.name}`,
  }));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!name.trim()) {
      setError("Name is required");
      return;
    }
    if (!hostnames.trim()) {
      setError("At least one hostname is required");
      return;
    }
    if (!selectedGateway) {
      setError("Gateway selection is required");
      return;
    }
    if (createSecurityPolicy && !selectedSecret) {
      setError("Htpasswd Secret selection is required when SecurityPolicy is enabled");
      return;
    }

    setLoading(true);
    setError(null);

    try {
      // Parse hostnames (one per line)
      const hostnameList = hostnames
        .split("\n")
        .map((h) => h.trim())
        .filter((h) => h.length > 0);

      // Parse selected gateway
      const [ns, gwName] = selectedGateway.split("/");

      const data: HTTPRouteCreateRequest = {
        name: name.trim(),
        namespace: service.namespace,
        hostnames: hostnameList,
        serviceName: service.name,
        servicePort: port,
        parentRefs: [
          {
            name: gwName,
            namespace: ns,
          },
        ],
      };

      // Add securityPolicy if enabled
      if (createSecurityPolicy) {
        data.securityPolicy = {
          basicAuthSecretName: selectedSecret,
        };
      }

      await createHTTPRoute(data);
      onCreate();
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create HTTPRoute");
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      className="max-w-[500px] p-6"
    >
      <div className="mb-4">
        <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
          Create HTTPRoute
        </h3>
        <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
          Service: <span className="font-medium">{service.namespace}/{service.name}</span>
        </p>
      </div>

      <form onSubmit={handleSubmit}>
        {error && (
          <div className="mb-4 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/20 dark:text-red-400">
            {error}
          </div>
        )}

        <div className="space-y-4">
          {/* Name */}
          <div>
            <Label htmlFor="route-name">HTTPRoute Name</Label>
            <Input
              id="route-name"
              type="text"
              placeholder="my-route"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />
          </div>

          {/* Hostnames */}
          <div>
            <Label htmlFor="hostnames">Hostnames (one per line)</Label>
            <TextArea
              placeholder="example.com&#10;www.example.com"
              rows={3}
              value={hostnames}
              onChange={setHostnames}
            />
          </div>

          {/* Gateway Parent */}
          <div>
            <Label htmlFor="gateway">Parent Gateway</Label>
            <Select
              key={key}
              options={gatewayOptions}
              placeholder="Select a Gateway"
              onChange={setSelectedGateway}
            />
          </div>

          {/* Port */}
          <div>
            <Label htmlFor="port">Service Port</Label>
            <Input
              type="number"
              min="1"
              max="65535"
              value={port.toString()}
              onChange={(e) => setPort(parseInt(e.target.value) || 80)}
            />
          </div>

          {/* SecurityPolicy */}
          <div className="border-t border-gray-200 dark:border-gray-700 pt-4 mt-4">
            <div className="mb-3">
              <label className="flex items-center gap-2">
                <input
                  type="checkbox"
                  id="create-security-policy"
                  checked={createSecurityPolicy}
                  onChange={(e) => setCreateSecurityPolicy(e.target.checked)}
                  className="h-4 w-4 rounded border-gray-300 text-brand-500 focus:ring-brand-500"
                />
                <Label htmlFor="create-security-policy" className="mb-0">
                  Also create SecurityPolicy with BasicAuth
                </Label>
              </label>
            </div>

            {createSecurityPolicy && (
              <div>
                <Label htmlFor="htpasswd-secret">Htpasswd Secret</Label>
                <Select
                  options={htpasswdSecrets.map((s) => ({
                    value: s.name,
                    label: s.name,
                  }))}
                  placeholder="Select an Htpasswd Secret"
                  onChange={setSelectedSecret}
                />
                <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                  Select an existing Htpasswd Secret to use for BasicAuth
                </p>
              </div>
            )}
          </div>
        </div>

        <div className="flex items-center justify-end gap-3 mt-6">
          <Button
            size="sm"
            variant="outline"
            onClick={onClose}
            disabled={loading}
          >
            Cancel
          </Button>
          <Button size="sm" disabled={loading}>
            {loading ? "Creating..." : "Create HTTPRoute"}
          </Button>
        </div>
      </form>
    </Modal>
  );
}
