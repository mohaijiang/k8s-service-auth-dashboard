"use client";

import React, { useState, useEffect, useCallback } from "react";
import { Modal } from "@/components/ui/modal";
import Button from "@/components/ui/button/Button";
import ComponentCard from "@/components/common/ComponentCard";
import HTTPRouteTable from "@/components/tables/HTTPRouteTable";
import HTTPRouteCreateModal from "@/components/modals/HTTPRouteCreateModal";
import type { ServiceOverview, HTTPRouteDetail, GatewayData } from "@/lib/api";
import {
  listHTTPRoutesByService,
  listGateways,
  deleteHTTPRoute,
} from "@/lib/api";

interface HTTPRouteManagerProps {
  isOpen: boolean;
  onClose: () => void;
  service: ServiceOverview;
}

export default function HTTPRouteManager({
  isOpen,
  onClose,
  service,
}: HTTPRouteManagerProps) {
  const [routes, setRoutes] = useState<HTTPRouteDetail[]>([]);
  const [gateways, setGateways] = useState<GatewayData[]>([]);
  const [loading, setLoading] = useState(false);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async () => {
    if (!isOpen) return;

    setLoading(true);
    setError(null);

    try {
      // Fetch routes for this service
      const routesRes = await listHTTPRoutesByService(
        service.namespace,
        service.name
      );
      setRoutes(routesRes.data);

      // Fetch available gateways
      const gatewaysRes = await listGateways();
      setGateways(gatewaysRes.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load data");
    } finally {
      setLoading(false);
    }
  }, [isOpen, service]);

  useEffect(() => {
    fetchData();
  }, [fetchData]);

  const handleCreate = () => {
    // Refresh data after creation
    fetchData();
  };

  const handleDelete = async (namespace: string, name: string) => {
    if (!confirm(`Are you sure you want to delete HTTPRoute ${name}?`)) {
      return;
    }

    try {
      await deleteHTTPRoute(namespace, name);
      // Refresh data after deletion
      fetchData();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete HTTPRoute");
    }
  };

  return (
    <>
      <Modal isOpen={isOpen} onClose={onClose} className="max-w-[800px] p-6">
        <div className="mb-4">
          <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
            Manage HTTPRoutes
          </h3>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            Service: <span className="font-medium">{service.namespace}/{service.name}</span>
          </p>
        </div>

        {error && (
          <div className="mb-4 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700 dark:border-red-800 dark:bg-red-900/20 dark:text-red-400">
            {error}
          </div>
        )}

        <ComponentCard title="HTTPRoutes" desc="Routing rules for this service">
          <div className="mb-4 flex justify-end">
            <Button size="sm" onClick={() => setShowCreateModal(true)}>
              Create HTTPRoute
            </Button>
          </div>

          <HTTPRouteTable
            routes={routes}
            loading={loading}
            onDelete={handleDelete}
          />
        </ComponentCard>

        <div className="flex items-center justify-end gap-3 mt-6">
          <Button size="sm" variant="outline" onClick={onClose}>
            Close
          </Button>
        </div>
      </Modal>

      {showCreateModal && (
        <HTTPRouteCreateModal
          isOpen={showCreateModal}
          onClose={() => setShowCreateModal(false)}
          service={service}
          gateways={gateways}
          onCreate={handleCreate}
        />
      )}
    </>
  );
}
